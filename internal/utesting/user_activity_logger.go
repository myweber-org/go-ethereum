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

func (l *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	
	l.handler.ServeHTTP(lw, r)
	
	duration := time.Since(start)
	log.Printf("%s %s %d %s", r.Method, r.URL.Path, lw.statusCode, duration)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}