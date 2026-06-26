package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hugr-lab/airport-go"
	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/central-set-go/internal/flight"
	"google.golang.org/grpc/credentials"

	// TELEMETRY
	// AUTOCERT
	"golang.org/x/crypto/acme/autocert"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = -10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func (app *application) loadTLSConfig() (*tls.Config, error) {
	enableTLS := strings.ToLower(os.Getenv("ENABLE_TLS")) == "true"
	if enableTLS {
		certFile := os.Getenv("TLS_CERT_FILE")
		keyFile := os.Getenv("TLS_KEY_FILE")
		caFile := os.Getenv("TLS_CA_CERT_FILE")
		if certFile == "" || keyFile == "" || caFile == "" {
			return nil, fmt.Errorf("ENABLE_TLS is true but TLS_CERT_FILE or TLS_KEY_FILE or TLS_CA_CERT_FILE is not set %s", "")
		}
		serverCert, err := tls.LoadX509KeyPair(
			certFile,
			keyFile,
		)
		if err != nil {
			return nil, fmt.Errorf("load server cert: %w", err)
		}
		// CA pool (for mTLS or client verification)
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read CA cert: %w", err)
		}
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("append CA cert")
		}
		return &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			// 👇 change this depending on your needs
			ClientAuth: tls.VerifyClientCertIfGiven,
			ClientCAs:  certPool,
			MinVersion: tls.VersionTLS12,
			// Good hygiene
			PreferServerCipherSuites: true,
		}, nil
	}
	return nil, nil
}

func (app *application) loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server certificate and key
	enableTLS := strings.ToLower(os.Getenv("ENABLE_TLS")) == "true"
	if enableTLS {
		certFile := os.Getenv("TLS_CERT_FILE")
		keyFile := os.Getenv("TLS_KEY_FILE")
		caFile := os.Getenv("TLS_CA_CERT_FILE")
		if certFile == "" || keyFile == "" {
			return nil, fmt.Errorf("ENABLE_TLS is true but TLS_CERT_FILE or TLS_KEY_FILE or TLS_CA_CERT_FILE is not set %s", "")
		}
		serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load server cert: %w", err)
		}
		// Load CA certificate for mutual TLS (optional)
		//if caFile != "" {
		certPool := x509.NewCertPool()
		if caCert, err := os.ReadFile(caFile); err == nil {
			if !certPool.AppendCertsFromPEM(caCert) {
				return nil, fmt.Errorf("failed to add CA cert to pool")
			}
		}
		//}
		// Configure TLS
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.NoClientCert, // Change to tls.RequireAndVerifyClientCert for mTLS
			ClientCAs:    certPool,
			MinVersion:   tls.VersionTLS12,
		}
		return credentials.NewTLS(tlsConfig), nil
	}
	return nil, nil
}

func (app *application) serveHTTP() error {
	//tlsCredentials, err := app.loadTLSCredentials()
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if env.GetBool("OTEL_ENABLED", false) {
		/*/ Set up OpenTelemetry.
		otelShutdown, err := setupOTelSDK(ctx)
		if err != nil {
			return err
		}
		// Handle shutdown properly so nothing leaks.
		defer func() {
			err = errors.Join(err, otelShutdown(context.Background()))
		}()*/
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.httpPort),
		BaseContext:  func(net.Listener) context.Context { return ctx },
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		//TLSConfig: ,
	}
	enableTLS := env.GetBool("ENABLE_TLS", false)
	autoCert := env.GetBool("AUTO_CERT", false)
	shutdownErrorChan := make(chan error)
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		/*ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()*/

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()
	app.logger.Info("starting server", slog.Group("server", "addr", srv.Addr))
	if enableTLS {
		if autoCert {
			domain := env.GetString("DOMAIN", "")
			if domain == "" {
				return fmt.Errorf("DOMAIN is required when AUTO_CERT=true")
			}
			autocertCache := env.GetString("AUTO_CERT_CACHE", "./certs")
			err := os.MkdirAll(autocertCache, 0755)
			if err != nil {
				return err
			}
			certManager := &autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				Cache:      autocert.DirCache(autocertCache),
				Email:      env.GetString("AUTO_CERT_EMAIL", ""),
				HostPolicy: autocert.HostWhitelist(domain),
			}
			//srv.TLSConfig = certManager.TLSConfig()
			srv.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}
			// HTTP challenge server
			/*go func() {
				challengeServer := &http.Server{
					Addr:    ":80",
					Handler: certManager.HTTPHandler(nil),
				}
				app.logger.Info("🔐 HTTPS server listening on with autocert", srv.Addr)
				if err := challengeServer.ListenAndServe(); err != nil &&
					!errors.Is(err, http.ErrServerClosed) {
					app.logger.Error(err.Error())
				}
			}()*/
			go func() {
				err = http.ListenAndServe(":80", certManager.HTTPHandler(nil))
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					app.logger.Error(err.Error())
				}
			}()
			err = srv.ListenAndServeTLS("", "")
			app.logger.Info("🔐 HTTPS server listening on with autocert", srv.Addr)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				app.logger.Error(err.Error())
			}
		} else {
			tlsConfig, err := app.loadTLSConfig()
			if err != nil && enableTLS {
				return err
			}
			srv.TLSConfig = tlsConfig
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				log.Fatal(err)
			}
			tlsListener := tls.NewListener(ln, srv.TLSConfig)
			app.logger.Info("🔐 HTTPS server listening on", srv.Addr)
			err = srv.Serve(tlsListener)
			if !errors.Is(err, http.ErrServerClosed) {
				return err
			}
		}
	} else {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	err := <-shutdownErrorChan
	if err != nil {
		return err
	}
	app.logger.Info("stopped server", slog.Group("server", "addr", srv.Addr))
	app.wg.Wait()
	return nil
}

func (app *application) serveSSE() error {
	//tlsCredentials, err := app.loadTLSCredentials()
	app.SSE_Broker = NewBroker()
	mux := http.NewServeMux()
	// SSE and notify endpoints
	mux.HandleFunc("/events", app.SSE_Broker.SSEHandler)
	mux.HandleFunc("/notify", app.SSE_Broker.NotifyHandler)
	// WS
	app.WS_ConnectionManager = app.NewConnectionManager()
	mux.HandleFunc("/ws", app.websocketEndpoint(app.WS_ConnectionManager))
	// SERVER
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.GetInt("SSE_SERVER_PORT", 5555)),
		Handler:      app.cors(mux),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
	enableTLS := strings.ToLower(os.Getenv("ENABLE_TLS")) == "true"
	tlsConfig, err := app.loadTLSConfig()
	if err != nil && enableTLS {
		return err
	}
	shutdownErrorChan := make(chan error)
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()
		shutdownErrorChan <- srv.Shutdown(ctx)
	}()
	app.logger.Info("starting SSE server", slog.Group("server", "addr", srv.Addr))
	if enableTLS {
		srv.TLSConfig = tlsConfig
		ln, err := net.Listen("tcp", srv.Addr)
		if err != nil {
			log.Fatal(err)
		}
		tlsListener := tls.NewListener(ln, srv.TLSConfig)
		app.logger.Info("🔐 HTTPS SSE server listening on", srv.Addr)
		err = srv.Serve(tlsListener)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	} else {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	err = <-shutdownErrorChan
	if err != nil {
		return err
	}
	app.logger.Info("stopped SSE server", slog.Group("server", "addr", srv.Addr))
	app.wg.Wait()
	return nil
}

func (app *application) airportValidateToken(token string) (string, error) {
	// fmt.Println("TOKEN:", token)
	identity, err := app.verifyTokenString(token)
	if err != nil {
		fmt.Printf("Err validating token: %s\n", err)
		return "", airport.ErrUnauthorized
	}
	jsonData, err := json.Marshal(identity)
	if err != nil {
		fmt.Printf("Err marshalling identity: %s\n", err)
		return "", airport.ErrUnauthorized
	}
	return string(jsonData), nil
}
func (app *application) validateToken(token string) (Dict, error) {
	identity, err := app.verifyTokenString(token)
	if err != nil {
		fmt.Printf("Err validating token: %s\n", err)
		return nil, err
	}
	return identity, nil
}
func (app *application) serveArrowFlight() error {
	// use crud.read to respect access control
	_sql := `SELECT * FROM "flight_catalog" WHERE active = ? AND excluded = ?`
	catalogs, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	_sql = `SELECT * FROM "flight_schema" WHERE active = ? AND excluded = ?`
	flight_schemas, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	_sql = `SELECT * FROM "flight_schema_table" WHERE active = ? AND excluded = ?`
	fligths_tables, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	_sql = `SELECT * FROM "flight_schema_table_field" WHERE active = ? AND excluded = ?`
	fligths_tables_fields, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	// flight_schema_table_scope
	_sql = `SELECT * FROM "flight_schema_table_scope" WHERE active = ? AND excluded = ?`
	fligths_tables_scopes, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	rla_tables := app.row_level_tables(Dict{"app": Dict{"app_id": 1, "db": app.config.db.dsn}})
	//fmt.Println(rla_tables)
	if !rla_tables["success"].(bool) {
		return fmt.Errorf("Error listing tables that requires RLA: %s", rla_tables["msg"])
	}
	// add tree tables and fields to flights, but with map[string]any like map[string]any{f["table"]: {"fields": map[string]any{"field": field}}}
	for _, f := range flight_schemas {
		var tables Dict = make(Dict)
		for _, t := range fligths_tables {
			if t["flight_schema_id"] == f["flight_schema_id"] {
				var fields Dict = make(Dict)
				for _, tf := range fligths_tables_fields {
					if tf["flight_schema_table_id"] == t["flight_schema_table_id"] {
						fields[tf["flight_schema_table_field"].(string)] = tf
					}
				}
				t["fields"] = fields
				var scopes Dict = make(Dict)
				for _, ts := range fligths_tables_scopes {
					if ts["flight_schema_table_id"] == t["flight_schema_table_id"] {
						scopes[ts["flight_schema_table_scope"].(string)] = ts
					}
				}
				t["scopes"] = scopes
				tables[t["flight_schema_table"].(string)] = t
			}
		}
		f["tables"] = tables
		f["rla_tables"] = rla_tables["tables"]
		if _conf, ok := f["flight_schema_conf"].(string); ok {
			var conf map[string]any
			err := json.Unmarshal([]byte(_conf), &conf)
			if err != nil {
				fmt.Printf("failed to parse flight_schema_conf JSON: %v\n", err)
			}
			f["conf"] = conf
		}
		//fmt.Println(f["rla_tables"])
	}
	fmt.Printf("#C: %d #F: %d\n", len(catalogs), len(flight_schemas))
	for _, c := range catalogs {
		schemas := []Dict{}
		for _, f := range flight_schemas {
			fmt.Printf("CID: %d FID: %d\n", c["flight_catalog_id"], f["flight_catalog_id"])
			if f["flight_catalog_id"] == c["flight_catalog_id"] {
				schemas = append(schemas, f)
			}
		}
		c["schemas"] = schemas
	}
	// start server
	// Create Flight adapter (airport-go) backed by our manager.
	//flightMgr := flight.NewAirportAdapter(flight_schemas, app.airportValidateToken, app.table_access, app.row_level_access, app.read)
	flightMgr := flight.NewAirportMultiCatalogsAdapter(catalogs, app.airportValidateToken, app.table_access, app.row_level_access, app.read)
	addr := env.GetString("ARROW_FLIGHT_ADDR", "0.0.0.0:50051")
	// Start the server (includes starting airport-go Flight server)
	if err := flightMgr.Start(addr); err != nil {
		return err
	}
	fmt.Printf("Arrow Flight Server Started @ %s", addr)
	// Wait for signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Print("shutting down...")
	// give Flight manager a chance to stop
	if err := flightMgr.Stop(context.Background()); err != nil {
		log.Printf("error during shutdown: %v", err)
	}
	// ensure manager shutdown (deferred previously) -- allow small wait
	_ = context.Background()
	log.Println("goodbye")
	return nil
}
