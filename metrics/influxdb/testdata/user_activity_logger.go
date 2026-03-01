package main

import (
    "encoding/json"
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

func logActivity(userID, action, details string) error {
    activity := ActivityLog{
        Timestamp: time.Now().UTC(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(activity); err != nil {
        return err
    }

    return nil
}

func main() {
    err := logActivity("user123", "login", "Successful authentication")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    err = logActivity("user456", "file_upload", "Uploaded profile picture")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    log.Println("Activity logging completed")
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp  time.Time
	Method     string
	Path       string
	RemoteAddr string
	UserAgent  string
	StatusCode int
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(lw, r)
		
		activity := ActivityLog{
			Timestamp:  start,
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			StatusCode: lw.statusCode,
		}
		
		log.Printf("ACTIVITY: %s %s %d %s %s %v",
			activity.Method,
			activity.Path,
			activity.StatusCode,
			activity.RemoteAddr,
			activity.UserAgent,
			activity.Timestamp.Format(time.RFC3339))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
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
	startTime := time.Now()
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(startTime)
	log.Printf(
		"Method: %s | Path: %s | Status: %d | Duration: %v | User-Agent: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.UserAgent(),
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