package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hugr-lab/airport-go"
	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/central-set-go/internal/flight"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = -10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func (app *application) serveHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.httpPort),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
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

	app.logger.Info("starting server", slog.Group("server", "addr", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", slog.Group("server", "addr", srv.Addr))

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
	_sql := `SELECT * FROM "arrow_flight" WHERE active = ? AND excluded = ?`
	fligths, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	_sql = `SELECT * FROM "arrow_flight_table" WHERE active = ? AND excluded = ?`
	fligths_tables, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	_sql = `SELECT * FROM "arrow_flight_table_field" WHERE active = ? AND excluded = ?`
	fligths_tables_fields, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	// add tree tables and fields to flights, but with map[string]any like map[string]any{f["table"]: {"fields": map[string]any{"field": field}}}
	for _, f := range fligths {
		var tables Dict = make(Dict)
		for _, t := range fligths_tables {
			if t["arrow_flight_id"] == f["arrow_flight_id"] {
				var fields Dict = make(Dict)
				for _, tf := range fligths_tables_fields {
					if tf["arrow_flight_table_id"] == t["arrow_flight_table_id"] {
						fields[tf["arrow_flight_table_field"].(string)] = tf
					}
				}
				t["fields"] = fields
				tables[t["arrow_flight_table"].(string)] = t
			}
		}
		f["tables"] = tables
	}
	// start server
	// Create Flight adapter (airport-go) backed by our manager.
	flightMgr := flight.NewAirportAdapter(fligths, app.airportValidateToken)
	addr := env.GetString("ARROW_FLIGHT_ADDR", "0.0.0.0:50051")
	// Start the server (includes starting airport-go Flight server)
	if err := flightMgr.Start(addr); err != nil {
		return err
	}
	fmt.Printf("server started at %s", addr)
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

func (app *application) serveArrowFlightTmp() error {
	// use crud.read to respect access control
	_sql := `SELECT * FROM "arrow_flight" WHERE active = ? AND excluded = ?`
	fligths, err := app.AdminGetRowsByFilter(_sql, []any{true, false})
	if err != nil {
		return err
	}
	//fmt.Println(_sql, fligths)
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return err
	}
	defer func() { // EXECUTE SHUTDOWN SQL ON EXIT
		fmt.Println("SHUTING DOWN FLIGHT SERVER")
		for _, f := range fligths { // execute shutdown_sql
			if _, ok := f["shutdown_sql"].(string); !ok {
				continue
			}
			fmt.Printf("%s: %s\n", f["arrow_flight"], f["shutdown_sql"])
			_, err := db.ExecContext(context.Background(), f["shutdown_sql"].(string))
			if err != nil {
				fmt.Printf("%s: %s: %s\n", f["arrow_flight"], f["shutdown_sql"], err)
			}
		}
		db.Close()
	}()
	// EXECUTE STARTUP SQL ON START
	fmt.Println("STARTING UP FLIGHT SERVER")
	for _, f := range fligths { // execute startup_sql
		if _, ok := f["startup_sql"].(string); !ok {
			continue
		}
		fmt.Printf("%s: %s\n", f["arrow_flight"], f["startup_sql"])
		_, err := db.ExecContext(context.Background(), f["startup_sql"].(string))
		if err != nil {
			fmt.Printf("%s: %s: %s\n", f["arrow_flight"], f["startup_sql"], err)
		}
		fmt.Printf("%s: %s\n", f["arrow_flight"], f["main_sql"])
		_, err = db.ExecContext(context.Background(), f["main_sql"].(string))
		if err != nil {
			fmt.Printf("%s: %s: %s\n", f["arrow_flight"], f["main_sql"], err)
		}
	}
	// start server
	// Create Flight adapter (airport-go) backed by our manager.
	flightMgr := flight.NewAirportAdapterTmp(db, fligths, app.airportValidateToken)

	addr := env.GetString("ARROW_FLIGHT_ADDR", "0.0.0.0:50051")
	// Start the server (includes starting airport-go Flight server)
	if err := flightMgr.Start(addr); err != nil {
		return err
	}
	fmt.Printf("server started at %s", addr)

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

// air
