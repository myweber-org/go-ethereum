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
		"Method: %s | Path: %s | Duration: %v | Timestamp: %s",
		r.Method,
		r.URL.Path,
		duration,
		time.Now().Format(time.RFC3339),
	)
}package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    IPAddress string    `json:"ip_address"`
}

type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.RWMutex
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

func (rl *RateLimiter) Allow(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    windowStart := now.Add(-rl.window)

    timestamps := rl.requests[ip]
    validRequests := make([]time.Time, 0)

    for _, ts := range timestamps {
        if ts.After(windowStart) {
            validRequests = append(validRequests, ts)
        }
    }

    if len(validRequests) >= rl.limit {
        return false
    }

    validRequests = append(validRequests, now)
    rl.requests[ip] = validRequests

    return true
}

func activityLogger(next http.Handler) http.Handler {
    limiter := NewRateLimiter(10, time.Minute)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr

        if !limiter.Allow(ip) {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }

        userID := r.Header.Get("X-User-ID")
        if userID == "" {
            userID = "anonymous"
        }

        logEntry := ActivityLog{
            Timestamp: time.Now().UTC(),
            UserID:    userID,
            Action:    r.Method + " " + r.URL.Path,
            IPAddress: ip,
        }

        logData, err := json.Marshal(logEntry)
        if err != nil {
            log.Printf("Failed to marshal log entry: %v", err)
        } else {
            log.Printf("Activity: %s", string(logData))
        }

        next.ServeHTTP(w, r)
    })
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    response := map[string]string{
        "status":  "success",
        "message": "Request processed",
    }
    json.NewEncoder(w).Encode(response)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", mainHandler)

    wrappedMux := activityLogger(mux)

    log.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
        log.Fatal(err)
    }
}