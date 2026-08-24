package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/realdatadriven/central-set-go/internal/env"
	"github.com/realdatadriven/central-set-go/internal/storage"
	// OPEN TELEMETRY
	//"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (app *application) S3Handler(w http.ResponseWriter, r *http.Request) {
	bucket := app.config.s3Bucket
	key := r.URL.Path[len("/uploads/"):]

	// Load AWS config (migrated awsConfig from before)
	cfg, err := app.awsConfig(context.Background())
	if err != nil {
		http.Error(w, "Failed to load AWS config", http.StatusInternalServerError)
		return
	}

	// Create S3 client
	svc := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint := app.config.s3Endpoint; endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		o.UsePathStyle = app.config.s3ForcePathStyle
	})

	// Get the object
	result, err := svc.GetObject(r.Context(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		http.Error(w, "Failed to get file from S3", http.StatusNotFound)
		return
	}
	defer result.Body.Close()

	// ContentType is now *string in v2 (still optional)
	if result.ContentType != nil {
		w.Header().Set("Content-Type", *result.ContentType)
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Stream file to client
	io.Copy(w, result.Body)
}

func (app *application) StorageAPIHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/uploads/")
	if key == "" {
		http.Error(w, "File not specified", http.StatusBadRequest)
		return
	}
	file, err := app.StorageAPI.Download(r.Context(), key)
	if err != nil {
		fmt.Println("KEY:", key)
		root := env.GetString("STORAGE_LOCAL_PATH", env.GetString("UPLOAD", "static/uploads"))
		strg, err2 := storage.NewLocalStorage(root)
		if err2 != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		} else {
			file, err2 = strg.Download(r.Context(), key)
			if err2 != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
		}
	}
	defer file.Close()
	// Detect content type from the file extension.
	contentType := mime.TypeByExtension(filepath.Ext(key))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	if _, err := io.Copy(w, file); err != nil {
		// At this point the response may already have started,
		// so there isn't much we can do with http.Error.
		return
	}
}

func (app *application) routes() http.Handler {
	//mux := httprouter.New()
	mux := http.NewServeMux()

	/*/ Register the WebSocket endpoint
	manager := app.NewConnectionManager()
	mux.HandleFunc("/ws", app.websocketEndpoint(manager))

	// Server-Sent Events (SSE)
	broker := NewBroker()
	mux.HandleFunc("/events", broker.SSEHandler)
	mux.HandleFunc("/sse", broker.SSEHandler)
	mux.HandleFunc("/notify", broker.NotifyHandler)*/

	// Handler for static files
	// Handler for static files with fallback to embedded assets
	// Tries filesystem first (for development/uploads), then embedded files
	fallbackServer := NewFallbackFileServer()
	mux.Handle("/static/", http.StripPrefix("/static/", fallbackServer))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fallbackServer))

	if app.StorageAPI != nil {
		mux.HandleFunc("/uploads/", app.StorageAPIHandler)
	} else if app.config.useS3 {
		mux.HandleFunc("/uploads/", app.S3Handler)
	} else {
		mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("static/uploads"))))
	}

	// Handler the root (index.html) - try filesystem first, then embedded
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		ServeStaticFile(w, r, "index.html")
	})

	// AI ASSISTANT ENDPOINTS
	/*{
		"provider": "googleai",
		"model": "gemini-1.5-flash",
		"system_prompt": "You are a ETLX assistant.",
		"system_prompt_file": "etlxllm.txt",
		"messages": [
				{"role": "user", "content": "Analyze this dataset"}
		]
	}*/
	//mux.HandleFunc("/etlx-assist", aiAssistHandler)

	//mux.NotFound = http.HandlerFunc(app.notFound)
	//mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.HandleFunc("GET /status", app.status)
	// mux.HandleFunc("POST /users", app.createUser)
	// mux.HandleFunc("POST /authentication-tokens", app.createAuthenticationToken)

	// REPLICA OF THE FASTAPI CENTRAL-SET
	//mux.HandleFunc("POST /dyn_api/login/login", app.dyn_api)
	mux.HandleFunc("POST /upload", app.uploadHandler)
	// OAuth2 token endpoint (minimal scaffold)
	mux.HandleFunc("POST /oauth2/token", app.oauthTokenHandler)
	// Protect dynamic API with OAuth middleware (validates JWT issued by existing login)
	mux.HandleFunc("POST /dyn_api/{ctrl}/{act}", app.dyn_api) //app.oauthMiddleware(app.dyn_api))

	// ODATA HANDDLER
	//mux.HandleFunc("GET /odata/{db}", app.odata_api_metadata)
	//mux.HandleFunc("GET /odata/{db}/", app.odata_api_metadata)
	//mux.HandleFunc("GET /odata/{db}", app.odata_api_service_document_xml)
	//mux.HandleFunc("GET /odata/{db}/", app.odata_api_service_document_xml)
	mux.HandleFunc("GET /odata/{db}", app.odata_api_service_document)
	mux.HandleFunc("GET /odata/{db}/", app.odata_api_service_document)
	mux.HandleFunc("GET /odata/{db}/$metadata", app.odata_api_metadata)
	mux.HandleFunc("GET /odata/{db}/$metadata/", app.odata_api_metadata)
	mux.HandleFunc("GET /odata/{db}/{table}", app.odata_api)
	mux.HandleFunc("GET /odata/{db}/{table}/", app.odata_api)
	mux.HandleFunc("GET /read_odata/{db}/{table}", app.read_odata)
	mux.HandleFunc("GET /read/{db}/{table}", app.read_odata)
	mux.HandleFunc("GET /rodata/{db}/{table}", app.read_odata)
	mux.HandleFunc("GET /odataodata/{db}/{table}", app.read_odata)

	// JOBS RUN ENDPOINTS
	mux.HandleFunc("GET /etlx/run/{name}", app.run_etlx_run_by_name)
	mux.HandleFunc("GET /etlx/name/{name}", app.run_etlx_run_by_name)
	mux.HandleFunc("GET /etlx/by_name/{name}", app.run_etlx_run_by_name)
	mux.HandleFunc("GET /etlx/run_by_name/{name}", app.run_etlx_run_by_name)
	mux.HandleFunc("GET /buckup", app.run_backup)
	mux.HandleFunc("GET /buckup/{name}", app.run_backup)
	mux.HandleFunc("GET /nb", app.run_notebook)
	mux.HandleFunc("GET /nb/{name}", app.run_notebook)
	mux.HandleFunc("GET /notebook", app.run_notebook)
	mux.HandleFunc("GET /notebook/{name}", app.run_notebook)
	mux.HandleFunc("GET /env", app.refreshEnv)
	mux.HandleFunc("GET /env/update", app.refreshEnv)
	mux.HandleFunc("GET /env/sync", app.refreshEnv)
	mux.HandleFunc("GET /env/refresh", app.refreshEnv)

	// QUACK ENDPOINTS
	mux.HandleFunc("GET /quack/start/{name}", app.startQuack)
	mux.HandleFunc("GET /quack/stop/{name}", app.stopQuack)
	mux.HandleFunc("GET /quack/restart/{name}", app.restartQuack)
	mux.HandleFunc("GET /quack_server/start/{name}", app.startQuack)
	mux.HandleFunc("GET /quack_server/stop/{name}", app.stopQuack)
	mux.HandleFunc("GET /quack_server/restart/{name}", app.restartQuack)
	mux.HandleFunc("GET /quack-server/start/{name}", app.startQuack)
	mux.HandleFunc("GET /quack-server/stop/{name}", app.stopQuack)
	mux.HandleFunc("GET /quack-server/restart/{name}", app.restartQuack)

	// OAUTH2
	mux.HandleFunc("GET /auth/{provider}/login", app.GothLoginHandler)
	mux.HandleFunc("GET /auth/{provider}/callback", app.GothCallbackHandler)

	// DYN UI
	mux.HandleFunc("POST /ui/{ui_slug}/login", app.serve_ui_login)
	mux.HandleFunc("POST /ui/{ui_slug}/login/", app.serve_ui_login)
	mux.HandleFunc("GET /ui/{ui_slug}/logout", app.logoutHandler)
	mux.HandleFunc("GET /ui/{ui_slug}/logout/", app.logoutHandler)
	mux.HandleFunc("POST /ui/{ui_slug}/logout", app.logoutHandler)
	mux.HandleFunc("POST /ui/{ui_slug}/logout/", app.logoutHandler)
	mux.HandleFunc("GET /ui/{ui_slug}", app.serve_ui_page)
	mux.HandleFunc("GET /ui/{ui_slug}/", app.serve_ui_page)
	mux.HandleFunc("GET /ui/{ui_slug}/{page_key}", app.serve_ui_page)
	mux.HandleFunc("GET /ui/{ui_slug}/{page_key}/", app.serve_ui_page)
	mux.HandleFunc("GET /ui/{ui_slug}/partial/{partial_name}", app.serve_ui_partial)
	mux.HandleFunc("GET /ui/{ui_slug}/partial/{partial_name}/", app.serve_ui_partial)
	mux.HandleFunc("GET /ui/{ui_slug}/static/{asset...}", app.serve_ui_asset)
	mux.HandleFunc("GET /ui/{ui_slug}/asset/{asset...}", app.serve_ui_asset)
	// OAUTH2
	mux.HandleFunc("GET /ui/{ui_slug}/oauth/{provider}/login", app.GothLoginHandler)
	mux.HandleFunc("GET /ui/{ui_slug}/oauth/{provider}/callback", app.HyperMGothCallbackHandler)
	// REST CRUD
	mux.HandleFunc("GET /crud/{db}/{table}", app.crud_api_handler)
	mux.HandleFunc("GET /crud/{db}/{table}/", app.crud_api_handler)
	mux.HandleFunc("GET /crud/{db}/{table}/{id}", app.crud_api_handler)
	// CRUD REST
	mux.HandleFunc("POST /crud/{db}/{table}", app.crud_api_handler)
	mux.HandleFunc("POST /crud/{db}/{table}/", app.crud_api_handler)
	mux.HandleFunc("PUT /crud/{db}/{table}/{id}", app.crud_api_handler)
	mux.HandleFunc("PATCH /crud/{db}/{table}/{id}", app.crud_api_handler)
	mux.HandleFunc("DELETE /crud/{db}/{table}/{id}", app.crud_api_handler)
	//http.HandleFunc("/ws", app.websocketEndpoint(manager))
	//app.rateLimit() || app.rateLimitMiddleware()
	/*/ OPEN TELEMETRY
	if env.GetBool("OTEL_ENABLED", false) {
		handler := otelhttp.NewHandler(mux, "/")
		return app.rateLimit(
			app.compress(
				app.cors(
					app.logAccess(
						app.recoverPanic(
							app.authenticate(handler),
						),
					),
				),
			),
		)
	} else {*/
	return app.rateLimit(
		app.compress(
			app.cors(
				app.logAccess(
					app.recoverPanic(
						app.sizeGuardMiddleware(
							app.authenticate(mux),
						),
					),
				),
			),
		),
	)
	//}
	//
}
