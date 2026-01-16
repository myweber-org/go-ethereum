package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	rateLimiter map[string][]time.Time
	window      time.Duration
	maxRequests int
}

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	UserAgent string    `json:"user_agent"`
	Status    int       `json:"status"`
}

func NewActivityLogger(window time.Duration, maxRequests int) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: make(map[string][]time.Time),
		window:      window,
		maxRequests: maxRequests,
	}
}

func (al *ActivityLogger) isRateLimited(ip string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	if requests, exists := al.rateLimiter[ip]; exists {
		var validRequests []time.Time
		for _, t := range requests {
			if t.After(windowStart) {
				validRequests = append(validRequests, t)
			}
		}
		al.rateLimiter[ip] = validRequests

		if len(validRequests) >= al.maxRequests {
			return true
		}
	}

	al.rateLimiter[ip] = append(al.rateLimiter[ip], now)
	return false
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if al.isRateLimited(ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)

		logEntry := ActivityLog{
			Timestamp: start,
			IP:        ip,
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			Status:    recorder.statusCode,
		}

		logData, _ := json.Marshal(logEntry)
		println(string(logData))
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}