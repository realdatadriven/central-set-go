package main

import (
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Handler serves files from an S3 bucket
func (app *application) S3Handler(w http.ResponseWriter, r *http.Request) {
	// Get the S3 bucket and key from the environment or request
	bucket := app.config.s3Bucket
	key := r.URL.Path[len("/uploads/"):]
	sess, err := app.awsSession()
	if err != nil {
		http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
		return
	}
	// Create an S3 service client
	svc := s3.New(sess)
	// Get the file from S3
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		http.Error(w, "Failed to get file from S3", http.StatusNotFound)
		return
	}
	defer result.Body.Close()
	// Set the correct content type and serve the file
	w.Header().Set("Content-Type", *result.ContentType)
	io.Copy(w, result.Body)
}

func (app *application) routes() http.Handler {
	//mux := httprouter.New()
	mux := http.NewServeMux()

	// Handler for static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets"))))
	if app.config.useS3 {
		mux.HandleFunc("/uploads/", app.S3Handler)
	} else {
		mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("static/uploads"))))
	}
	// Handler the root (index.html)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	//mux.NotFound = http.HandlerFunc(app.notFound)
	//mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.HandleFunc("GET /status", app.status)
	// mux.HandleFunc("POST /users", app.createUser)
	// mux.HandleFunc("POST /authentication-tokens", app.createAuthenticationToken)

	// REPLICA OF THE FASTAPI CENTRAL-SET
	mux.HandleFunc("POST /dyn_api/login/login", app.login)
	mux.HandleFunc("POST /upload", app.uploadHandler)
	mux.HandleFunc("POST /dyn_api/{ctrl}/{act}", app.dyn_api)
	// RUN ENDPOINTS
	mux.HandleFunc("GET /etlx/run/{name}", app.run_etlx_run_by_name)

	//mux.Handler("GET /protected", app.requireAuthenticatedUser(http.HandleFunc(app.protected)))

	//mux.Handler("GET /basic-auth-protected", app.requireBasicAuthentication(http.HandleFunc(app.protected)))

	// Register the WebSocket endpoint
	manager := app.NewConnectionManager()
	//mux.HandleFunc("/ws", app.websocketEndpoint(manager))
	http.HandleFunc("/ws", app.websocketEndpoint(manager))
	return app.compress(app.cors(app.logAccess(app.recoverPanic(app.authenticate(mux)))))
}
