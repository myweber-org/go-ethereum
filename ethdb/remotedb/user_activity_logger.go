package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	rateLimiter *RateLimiter
}

type RateLimiter struct {
	requests map[string][]time.Time
	interval time.Duration
	max      int
}

func NewRateLimiter(interval time.Duration, maxRequests int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		interval: interval,
		max:      maxRequests,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	timestamps := rl.requests[ip]

	var valid []time.Time
	for _, ts := range timestamps {
		if now.Sub(ts) <= rl.interval {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= rl.max {
		return false
	}

	valid = append(valid, now)
	rl.requests[ip] = valid
	return true
}

func NewActivityLogger() *ActivityLogger {
	return &ActivityLogger{
		rateLimiter: NewRateLimiter(time.Minute, 100),
	}
}

func (al *ActivityLogger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		
		if !al.rateLimiter.Allow(ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		start := time.Now()
		log.Printf("Activity started: %s %s from %s", r.Method, r.URL.Path, ip)

		defer func() {
			duration := time.Since(start)
			log.Printf("Activity completed: %s %s from %s took %v", 
				r.Method, r.URL.Path, ip, duration)
		}()

		next.ServeHTTP(w, r)
	})
}