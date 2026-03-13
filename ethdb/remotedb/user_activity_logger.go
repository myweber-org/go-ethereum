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
}package middleware

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
	recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
	
	al.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	
	log.Printf(
		"%s %s %d %s %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityLog struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Details   string    `json:"details,omitempty"`
}

func NewActivityLog(userID, action, details string) *ActivityLog {
    return &ActivityLog{
        Timestamp: time.Now().UTC(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }
}

func (al *ActivityLog) ToJSON() ([]byte, error) {
    return json.MarshalIndent(al, "", "  ")
}

func LogActivity(userID, action, details string) error {
    logEntry := NewActivityLog(userID, action, details)
    jsonData, err := logEntry.ToJSON()
    if err != nil {
        return fmt.Errorf("failed to marshal activity log: %w", err)
    }

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    if _, err := file.Write(append(jsonData, '\n')); err != nil {
        return fmt.Errorf("failed to write log entry: %w", err)
    }

    log.Printf("Activity logged: %s performed %s", userID, action)
    return nil
}

func main() {
    if err := LogActivity("user123", "login", "Successful authentication"); err != nil {
        log.Fatalf("Failed to log activity: %v", err)
    }

    if err := LogActivity("user456", "file_upload", "Uploaded document.pdf"); err != nil {
        log.Fatalf("Failed to log activity: %v", err)
    }

    fmt.Println("Activity logging completed")
}