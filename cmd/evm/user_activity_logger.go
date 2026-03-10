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
	
	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Path      string
	Method    string
	Timestamp time.Time
	IPAddress string
}

type ActivityLogger struct {
	activityChan chan ActivityLog
}

func NewActivityLogger(bufferSize int) *ActivityLogger {
	al := &ActivityLogger{
		activityChan: make(chan ActivityLog, bufferSize),
	}
	go al.processLogs()
	return al
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		activity := ActivityLog{
			UserID:    userID,
			Path:      r.URL.Path,
			Method:    r.Method,
			Timestamp: time.Now(),
			IPAddress: r.RemoteAddr,
		}

		select {
		case al.activityChan <- activity:
		default:
			log.Println("Activity log buffer full, dropping entry")
		}

		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) processLogs() {
	for activity := range al.activityChan {
		log.Printf("ACTIVITY: User=%s IP=%s %s %s at %s",
			activity.UserID,
			activity.IPAddress,
			activity.Method,
			activity.Path,
			activity.Timestamp.Format(time.RFC3339))
	}
}

func (al *ActivityLogger) Close() {
	close(al.activityChan)
}