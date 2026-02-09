package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/realdatadriven/central-set-go/internal/response"

	"github.com/pascaldekloe/jwt"
	"github.com/tomasen/realip"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) logAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := response.NewMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		var (
			ip     = realip.FromRequest(r)
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
		responseAttrs := slog.Group("repsonse", "status", mw.StatusCode, "size", mw.BytesCount)

		app.logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")

			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				token := headerParts[1]

				claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secretKey))
				if err != nil {
					app.invalidAuthenticationToken(w, r)
					return
				}

				if !claims.Valid(time.Now()) {
					app.expiredAuthenticationToken(w, r)
					return
				}

				if claims.Issuer != app.config.baseURL {
					app.invalidAuthenticationToken(w, r)
					return
				}

				if !claims.AcceptAudience(app.config.baseURL) {
					app.invalidAuthenticationToken(w, r)
					return
				}

				/*userID, err := strconv.Atoi(claims.Subject)
				if err != nil {
					app.serverError(w, r, err)
					return
				}

				user, found, err := app.db.GetUserById(userID)
				if err != nil {
					app.serverError(w, r, err)
					return
				}

				if found {
					r = contextSetAuthenticatedUser(r, user)
				}*/
				var user map[string]interface{}
				//print(1, " ", claims.Subject, "\n")
				err2 := json.Unmarshal([]byte(claims.Subject), &user)
				if err2 == nil {
					//print(2, " ", user["username"].(string), "\n")
					r = contextSetAuthenticatedUser(r, &user)
				}
			} else {
				//app.invalidAuthenticationToken(w, r)
				//return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := contextGetAuthenticatedUser(r)

		if authenticatedUser == nil {
			app.authenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireBasicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, plaintextPassword, ok := r.BasicAuth()
		if !ok {
			app.basicAuthenticationRequired(w, r)
			return
		}

		if app.config.basicAuth.username != username {
			app.basicAuthenticationRequired(w, r)
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(app.config.basicAuth.hashedPassword), []byte(plaintextPassword))
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			app.basicAuthenticationRequired(w, r)
			return
		case err != nil:
			app.serverError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, enctype")

		// Handle preflight (OPTIONS) request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)

}

func shouldCompress(r *http.Request) bool {
	// Verifica se o cliente aceita codificação gzip
	acceptEncoding := r.Header.Get("Accept-Encoding")
	return strings.Contains(acceptEncoding, "gzip")
}

func (app *application) compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shouldCompress(r) {
			next.ServeHTTP(w, r)
			return
		}
		// Criar um escritor gzip
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Substituir o escritor original por um que escreve para o escritor gzip
		rw := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(rw, r)
	})
}

// Rate Limiting: use inmemory duckdb to store request counts per IP and reset them after a certain time window. If an IP exceeds the limit, return a 429 Too Many Requests response. This can help prevent abuse and ensure fair usage of your API.

func (app *application) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.rateLimitingEnabled  {
			ip := realip.FromRequest(r)
			// Check if the IP has exceeded the rate limit
			if app.IsRateLimited(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
				return
			}
			// Increment the request count for the IP
			app.Increment(ip)
		}
		next.ServeHTTP(w, r)
	})
}

// create app.rateLimiter with a simple in-memory duckdb databse in app.memdb

func (app *application) IsRateLimited(ip string) bool {
	var requestCount int
	var lastRequestTime time.Time
	err := app.memdb.QueryRow("SELECT request_count, last_request_time FROM rate_limits WHERE ip = ?", ip).Scan(&requestCount, &lastRequestTime)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error checking rate limit:", err)
		return false
	}
	if time.Since(lastRequestTime) > time.Minute {
		return false
	}
	return requestCount >= app.rtRequestLimit // Example limit: 100 requests per minute
}

func (app *application) Increment(ip string) (int, error) {
    query := fmt.Sprintf(`MERGE INTO rate_limits
USING (SELECT '%s' as ip, 1 as request_count, CURRENT_TIMESTAMP as last_request_time) AS upserts
ON rate_limits.ip = upserts.ip
WHEN MATCHED THEN UPDATE SET 
    request_count = rate_limits.request_count + 1,
    last_request_time = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN INSERT VALUES (upserts.ip, upserts.request_count, upserts.last_request_time);`, ip)
	return rl.db.Exec(query)
}
