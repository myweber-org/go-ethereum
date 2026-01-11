
package middleware

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ActivityLog struct {
	UserID    string
	Action    string
	Timestamp time.Time
	IPAddress string
}

type ActivityLogger struct {
	limiter   *rate.Limiter
	eventChan chan<- ActivityLog
}

func NewActivityLogger(eventsPerSecond float64, burst int, eventChan chan<- ActivityLog) *ActivityLogger {
	return &ActivityLogger{
		limiter:   rate.NewLimiter(rate.Limit(eventsPerSecond), burst),
		eventChan: eventChan,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := extractUserID(r)
		action := r.Method + " " + r.URL.Path
		ip := r.RemoteAddr

		logEntry := ActivityLog{
			UserID:    userID,
			Action:    action,
			Timestamp: time.Now().UTC(),
			IPAddress: ip,
		}

		if al.limiter.Allow() {
			select {
			case al.eventChan <- logEntry:
			default:
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return auth[:min(8, len(auth))]
	}
	return "anonymous"
}