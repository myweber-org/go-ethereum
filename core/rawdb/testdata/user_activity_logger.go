package middleware

import (
	"log"
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

func NewActivityLogger(window time.Duration, maxRequests int) *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: make(map[string][]time.Time),
		window:      window,
		maxRequests: maxRequests,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		path := r.URL.Path
		method := r.Method

		if !al.allowRequest(clientIP) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("Activity: %s %s %s %s Duration: %v", clientIP, method, path, userAgent, duration)
	})
}

func (al *ActivityLogger) allowRequest(ip string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-al.window)

	if _, exists := al.rateLimiter[ip]; !exists {
		al.rateLimiter[ip] = []time.Time{}
	}

	validRequests := []time.Time{}
	for _, t := range al.rateLimiter[ip] {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= al.maxRequests {
		return false
	}

	validRequests = append(validRequests, now)
	al.rateLimiter[ip] = validRequests

	go al.cleanupOldEntries(ip, windowStart)

	return true
}

func (al *ActivityLogger) cleanupOldEntries(ip string, cutoff time.Time) {
	al.mu.Lock()
	defer al.mu.Unlock()

	if entries, exists := al.rateLimiter[ip]; exists {
		validEntries := []time.Time{}
		for _, t := range entries {
			if t.After(cutoff) {
				validEntries = append(validEntries, t)
			}
		}
		if len(validEntries) == 0 {
			delete(al.rateLimiter, ip)
		} else {
			al.rateLimiter[ip] = validEntries
		}
	}
}