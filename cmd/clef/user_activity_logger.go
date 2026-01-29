package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	IPAddress string
	Method    string
	Path      string
	Timestamp time.Time
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserIDFromToken(authHeader)
		}

		activity := ActivityLog{
			UserID:    userID,
			IPAddress: r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
			Timestamp: start,
		}

		log.Printf("Activity: %s %s by %s from %s", 
			activity.Method, 
			activity.Path, 
			activity.UserID, 
			activity.IPAddress)

		next.ServeHTTP(w, r)
	})
}

func extractUserIDFromToken(token string) string {
	return "user123"
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
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}