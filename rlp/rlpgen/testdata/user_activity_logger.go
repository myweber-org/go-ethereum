
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserID(authHeader)
		}

		activity := ActivityLog{
			Timestamp: start,
			UserID:    userID,
			Method:    r.Method,
			Path:      r.URL.Path,
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		log.Printf("Activity: %s %s by %s from %s", activity.Method, activity.Path, activity.UserID, activity.IPAddress)

		next.ServeHTTP(w, r)
	})
}

func extractUserID(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:min(len(authHeader), 15)]
	}
	return "authenticated"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}