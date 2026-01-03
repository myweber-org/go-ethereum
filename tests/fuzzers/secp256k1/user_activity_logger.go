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
		"Activity: %s %s | Status: %d | Duration: %v | UserAgent: %s | RemoteAddr: %s",
		r.Method,
		r.URL.Path,
		getStatusCode(w),
		duration,
		r.UserAgent(),
		r.RemoteAddr,
	)
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func getStatusCode(w http.ResponseWriter) int {
	if rw, ok := w.(*responseWriterWrapper); ok {
		return rw.statusCode
	}
	return http.StatusOK
}

func ActivityLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		logger := NewActivityLogger(next)
		logger.ServeHTTP(wrappedWriter, r)
	})
}