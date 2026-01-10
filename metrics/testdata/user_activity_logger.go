package middleware

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	activities  map[string][]ActivityRecord
	rateLimiter *RateLimiter
}

type ActivityRecord struct {
	UserID    string
	Action    string
	Timestamp time.Time
	IPAddress string
}

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func NewActivityLogger(limit int, window time.Duration) *ActivityLogger {
	return &ActivityLogger{
		activities: make(map[string][]ActivityRecord),
		rateLimiter: &RateLimiter{
			requests: make(map[string][]time.Time),
			limit:    limit,
			window:   window,
		},
	}
}

func (rl *RateLimiter) Allow(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	if _, exists := rl.requests[userID]; !exists {
		rl.requests[userID] = []time.Time{}
	}

	validRequests := []time.Time{}
	for _, t := range rl.requests[userID] {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) >= rl.limit {
		return false
	}

	validRequests = append(validRequests, now)
	rl.requests[userID] = validRequests
	return true
}

func (al *ActivityLogger) LogActivity(ctx context.Context, userID, action, ip string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if !al.rateLimiter.Allow(userID) {
			return ErrRateLimitExceeded
		}

		record := ActivityRecord{
			UserID:    userID,
			Action:    action,
			Timestamp: time.Now(),
			IPAddress: ip,
		}

		al.mu.Lock()
		al.activities[userID] = append(al.activities[userID], record)
		al.mu.Unlock()

		log.Printf("Activity logged: %s performed %s from %s", userID, action, ip)
		return nil
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err := al.LogActivity(ctx, userID, r.Method+" "+r.URL.Path, r.RemoteAddr)
		if err != nil {
			if err == ErrRateLimitExceeded {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			log.Printf("Failed to log activity: %v", err)
		}

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) GetUserActivities(userID string) []ActivityRecord {
	al.mu.RLock()
	defer al.mu.RUnlock()

	return append([]ActivityRecord{}, al.activities[userID]...)
}

var ErrRateLimitExceeded = error.New("rate limit exceeded")