// internal/auth/jwt.go
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims example – adapt to your needs
type CustomClaims struct {
	UserID    string `json:"sub,omitempty"`
	Provider  string `json:"provider,omitempty"` // google, github, company-okta, ...
	Email     string `json:"email,omitempty"`
	TokenType string `json:"type,omitempty"` // "access", "refresh"
	jwt.RegisteredClaims
}

// JWTSecretKey – in production use env var or secret manager
// Never hardcode!
var jwtSecret = []byte(os.Getenv("JWT_SIGNING_SECRET")) // must be 32+ bytes for HS256

// Must be very strong & rotated periodically
func init() {
	if len(jwtSecret) < 32 {
		panic("JWT_SIGNING_SECRET must be at least 32 bytes long")
	}
}

// ParseAndValidateToken returns validated claims or error
func ParseAndValidateToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Enforce algorithm – critical security step!
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	// Extra business rules (optional but recommended)
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// Middleware example (chi / gin / std http compatible)
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader { // no Bearer prefix
			http.Error(w, "invalid bearer token format", http.StatusUnauthorized)
			return
		}

		claims, err := ParseAndValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Attach claims to request context
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}