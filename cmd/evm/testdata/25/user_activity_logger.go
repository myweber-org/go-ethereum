package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	Logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{Logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		userAgent := r.UserAgent()
		clientIP := r.RemoteAddr
		method := r.Method
		path := r.URL.Path

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		status := recorder.statusCode

		al.Logger.Printf(
			"Activity: %s %s | Status: %d | Duration: %v | IP: %s | Agent: %s",
			method, path, status, duration, clientIP, userAgent,
		)
	})
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

var activityLogs []ActivityLog

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logEntry := ActivityLog{
            UserID:    extractUserID(r),
            Endpoint:  r.URL.Path,
            Method:    r.Method,
            Timestamp: time.Now(),
            IPAddress: r.RemoteAddr,
        }

        activityLogs = append(activityLogs, logEntry)
        log.Printf("Activity logged: %s %s by user %s", logEntry.Method, logEntry.Endpoint, logEntry.UserID)

        next.ServeHTTP(w, r)
    })
}

func extractUserID(r *http.Request) string {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return "anonymous"
    }
    return "user_" + authHeader[:8]
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Data endpoint response"))
    })

    loggedMux := loggingMiddleware(mux)

    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", loggedMux)
}