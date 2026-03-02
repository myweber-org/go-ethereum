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
	startTime := time.Now()
	writer := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	
	al.handler.ServeHTTP(writer, r)
	
	duration := time.Since(startTime)
	log.Printf("[%s] %s %s %d %v", 
		r.RemoteAddr, 
		r.Method, 
		r.URL.Path, 
		writer.statusCode, 
		duration)
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

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	UserAgent string
	IPAddress string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activity := ActivityLog{
			Timestamp: time.Now().UTC(),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
			IPAddress: r.RemoteAddr,
		}

		log.Printf("Activity: %s %s from %s (%s)", 
			activity.Method, 
			activity.Path, 
			activity.IPAddress, 
			activity.UserAgent)

		next.ServeHTTP(w, r)
	})
}