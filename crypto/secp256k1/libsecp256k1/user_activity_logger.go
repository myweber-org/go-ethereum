
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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

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
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Endpoint  string
	Method    string
	Timestamp time.Time
	IPAddress string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserID(authHeader)
		}

		activity := ActivityLog{
			UserID:    userID,
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: r.RemoteAddr,
		}

		log.Printf("Activity: %s %s by %s from %s", 
			activity.Method, 
			activity.Endpoint, 
			activity.UserID, 
			activity.IPAddress)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(token string) string {
	return "user_" + token[:8]
}package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

type ActivityLog struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

func logActivity(userID, action, details string) {
    activity := ActivityLog{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    logEntry, err := json.Marshal(activity)
    if err != nil {
        log.Printf("Failed to marshal activity: %v", err)
        return
    }

    fmt.Println(string(logEntry))

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Failed to open log file: %v", err)
        return
    }
    defer file.Close()

    if _, err := file.Write(append(logEntry, '\n')); err != nil {
        log.Printf("Failed to write to log file: %v", err)
    }
}

func main() {
    logActivity("user123", "login", "User logged in from web browser")
    logActivity("user456", "purchase", "Purchased item ID: 789")
    logActivity("user123", "logout", "Session terminated")
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
	
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
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
}