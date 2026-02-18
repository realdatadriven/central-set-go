package main

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ClientLimiter holds per-user rate limiters
type ClientLimiter struct {
	limiters map[string]*rate.Limiter // key = userID / apiKey / IP
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewClientLimiter creates a new per-client rate limiter manager
func NewClientLimiter(rps float64, burst int) *ClientLimiter {
	return &ClientLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// GetLimiter returns (or creates) limiter for a user
func (cl *ClientLimiter) GetLimiter(key string) *rate.Limiter {
	cl.mu.RLock()
	lim, exists := cl.limiters[key]
	cl.mu.RUnlock()
	if exists {
		return lim
	}
	// Upgrade to write lock
	cl.mu.Lock()
	defer cl.mu.Unlock()
	// Double-check after lock
	if lim, exists = cl.limiters[key]; exists {
		return lim
	}
	lim = rate.NewLimiter(cl.rate, cl.burst)
	cl.limiters[key] = lim
	return lim
}

// Optional: cleanup old entries (call periodically)
func (cl *ClientLimiter) Cleanup() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	for key, lim := range cl.limiters {
		// Example: remove if no burst capacity recently used
		if lim.Burst() == cl.burst && lim.AllowN(time.Now(), 0) {
			delete(cl.limiters, key)
		}
	}
}

var globalLimiter = NewClientLimiter(5, 20) // 5 req/s, burst 20 per user

func (app *application) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In real app: get user from auth header, JWT, session, etc.
		// Here we use IP as example
		if app.rateLimitingEnabled {
			userKey := r.RemoteAddr // or r.Header.Get("X-API-Key") or userID from context
			limiter := globalLimiter.GetLimiter(userKey)
			if !limiter.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
