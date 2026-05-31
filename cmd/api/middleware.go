package main

import (
	"compress/gzip"
	"database/sql"
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

		/*/ USE COOKIE PROVIDED IN THE OAUTH AS IF Authorization HEADER
		cookie, err := r.Cookie("session")
		if err == nil && authorizationHeader == "" {
			authorizationHeader = "Bearer " + cookie.Value
		}*/

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
		// ignore get static file req
		if app.rateLimitingEnabled && r.Method == "GET" && strings.Contains(r.URL.Path, "static") {
			next.ServeHTTP(w, r)
		}
		if app.rateLimitingEnabled {
			ip := realip.FromRequest(r)
			// Check if the IP has exceeded the rate limit
			if app.IsRateLimited(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": "Too many requests"})
				return
			}
			// Increment the request count for the IP
			_, err := app.Increment(ip)
			if err != nil {
				fmt.Println("Error incrementing rate limit:", err)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// create app.rateLimiter with a simple in-memory duckdb databse in app.memdb

func (app *application) IsRateLimited(ip string) bool {
	var requestCount int
	var lastRequestTime time.Time
	err := app.memdb.QueryRow("SELECT request_count, last_request_time FROM rate_limits WHERE ip = ? LIMIT 1", ip).Scan(&requestCount, &lastRequestTime)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Error checking rate limit:", err)
		return false
	} else if err == sql.ErrNoRows {
		lastRequestTime = time.Now()
		query := fmt.Sprintf(`INSERT INTO rate_limits VALUES ('%s', 1, CURRENT_TIMESTAMP);`, ip)
		_, err := app.memdb.Exec(query)
		if err != nil {
			fmt.Println("Error resetting rate limit:", err, query)
		}
	}
	//fmt.Println("RATE Limiting lastRequestTime:", requestCount, lastRequestTime, time.Since(lastRequestTime), time.Minute, time.Since(lastRequestTime) > time.Minute)
	//fmt.Println("RATE Limiting lastRequestTime CHECKING:", requestCount, lastRequestTime, time.Since(time.Now()), time.Minute, time.Since(time.Now()) > time.Minute)
	if time.Since(lastRequestTime) > time.Minute {
		// reset request count after 1 minute
		query := fmt.Sprintf(`MERGE INTO rate_limits
USING (SELECT '%s' as ip, 1 as request_count, CURRENT_TIMESTAMP as last_request_time) AS upserts
ON rate_limits.ip = upserts.ip
WHEN MATCHED THEN UPDATE SET request_count = 0, last_request_time = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN INSERT VALUES (upserts.ip, upserts.request_count, upserts.last_request_time);`, ip)
		// reset request count after 1 minute
		query = fmt.Sprintf(`UPDATE rate_limits SET request_count = 0, last_request_time = CURRENT_TIMESTAMP WHERE ip='%s';`, ip)
		_, err := app.memdb.Exec(query)
		if err != nil {
			fmt.Println("Error resetting rate limit:", err, query)
		}
		return false
	}
	//fmt.Println("REQUEST vs LIMIT:", requestCount, app.rtRequestLimit, requestCount >= app.rtRequestLimit)
	return requestCount >= app.rtRequestLimit // Example limit: 100 requests per minute
}

func (app *application) Increment(ip string) (sql.Result, error) {
	query := fmt.Sprintf(`MERGE INTO rate_limits
USING (SELECT '%s' as ip, 1 as request_count, CURRENT_TIMESTAMP as last_request_time) AS upserts
ON rate_limits.ip = upserts.ip
WHEN MATCHED THEN UPDATE SET 
    request_count = rate_limits.request_count + 1,
    last_request_time = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN INSERT VALUES (upserts.ip, upserts.request_count, upserts.last_request_time);`, ip)
	query = fmt.Sprintf(`UPDATE rate_limits SET request_count = request_count + 1, last_request_time = CURRENT_TIMESTAMP WHERE ip='%s';`, ip)
	return app.memdb.Exec(query)
}
