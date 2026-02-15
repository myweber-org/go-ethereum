package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Endpoint  string    `json:"endpoint"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

type ActivityLogger struct {
	rateLimit time.Duration
	lastLog   map[string]time.Time
}

func NewActivityLogger(rateLimit time.Duration) *ActivityLogger {
	return &ActivityLogger{
		rateLimit: rateLimit,
		lastLog:   make(map[string]time.Time),
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		key := fmt.Sprintf("%s-%s", userID, r.URL.Path)
		now := time.Now()

		if last, exists := al.lastLog[key]; exists {
			if now.Sub(last) < al.rateLimit {
				next.ServeHTTP(w, r)
				return
			}
		}

		activity := ActivityLog{
			Timestamp: now,
			UserID:    userID,
			Action:    r.Method,
			Endpoint:  r.URL.Path,
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		logData, err := json.Marshal(activity)
		if err == nil {
			fmt.Printf("ACTIVITY: %s\n", string(logData))
		}

		al.lastLog[key] = now
		next.ServeHTTP(w, r)
	})
}

func (al *ActivityLogger) CleanupOldLogs() {
	go func() {
		for {
			time.Sleep(time.Hour)
			now := time.Now()
			for key, lastTime := range al.lastLog {
				if now.Sub(lastTime) > time.Hour*24 {
					delete(al.lastLog, key)
				}
			}
		}
	}()
}