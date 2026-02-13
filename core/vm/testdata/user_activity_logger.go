package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf(
		"Activity: %s %s from %s completed in %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	userHits    map[string][]time.Time
	windowSize  time.Duration
	maxRequests int
}

func NewActivityLogger(window time.Duration, limit int) *ActivityLogger {
	return &ActivityLogger{
		userHits:    make(map[string][]time.Time),
		windowSize:  window,
		maxRequests: limit,
	}
}

func (al *ActivityLogger) cleanupOldHits(userID string) {
	now := time.Now()
	al.mu.Lock()
	defer al.mu.Unlock()

	hits := al.userHits[userID]
	validHits := []time.Time{}
	for _, hit := range hits {
		if now.Sub(hit) <= al.windowSize {
			validHits = append(validHits, hit)
		}
	}
	al.userHits[userID] = validHits
}

func (al *ActivityLogger) isRateLimited(userID string) bool {
	al.cleanupOldHits(userID)

	al.mu.RLock()
	hits := al.userHits[userID]
	al.mu.RUnlock()

	return len(hits) >= al.maxRequests
}

func (al *ActivityLogger) recordActivity(userID string) {
	al.mu.Lock()
	defer al.mu.Unlock()
	al.userHits[userID] = append(al.userHits[userID], time.Now())
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		if al.isRateLimited(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for user: %s", userID)
			return
		}

		al.recordActivity(userID)

		log.Printf("Activity: %s %s from user: %s", r.Method, r.URL.Path, userID)
		next.ServeHTTP(w, r)
	})
}