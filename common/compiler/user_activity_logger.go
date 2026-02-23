package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ActivityLog struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address"`
}

type RateLimiter struct {
	mu       sync.Mutex
	counters map[string]int
	window   time.Duration
	limit    int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		counters: make(map[string]int),
		window:   window,
		limit:    limit,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	count, exists := rl.counters[key]
	if !exists {
		go rl.resetCounter(key)
	}

	if count >= rl.limit {
		return false
	}

	rl.counters[key] = count + 1
	return true
}

func (rl *RateLimiter) resetCounter(key string) {
	time.Sleep(rl.window)
	rl.mu.Lock()
	delete(rl.counters, key)
	rl.mu.Unlock()
}

func ActivityLogger(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				userID = "anonymous"
			}

			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			key := userID + ":" + r.Method + ":" + r.URL.Path
			if !limiter.Allow(key) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			logEntry := ActivityLog{
				UserID:    userID,
				Action:    "request",
				Path:      r.URL.Path,
				Method:    r.Method,
				Timestamp: time.Now().UTC(),
				IPAddress: ip,
			}

			logData, err := json.Marshal(logEntry)
			if err == nil {
				go func(data []byte) {
					// In production, send to logging service
					println(string(data))
				}(logData)
			}

			next.ServeHTTP(w, r)
		})
	}
}