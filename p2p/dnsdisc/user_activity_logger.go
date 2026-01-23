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
	userID := extractUserID(r)
	ip := r.RemoteAddr

	al.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	log.Printf("User %s from IP %s accessed %s %s - Duration: %v",
		userID, ip, r.Method, r.URL.Path, duration)
}

func extractUserID(r *http.Request) string {
	if user := r.Header.Get("X-User-ID"); user != "" {
		return user
	}
	return "anonymous"
}