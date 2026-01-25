package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Activity struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip_address"`
}

type ActivityLogger struct {
	mu       sync.RWMutex
	activities []Activity
	rateLimit map[string]time.Time
}

func NewActivityLogger() *ActivityLogger {
	return &ActivityLogger{
		rateLimit: make(map[string]time.Time),
	}
}

func (al *ActivityLogger) LogActivity(userID, action, ip string) bool {
	al.mu.Lock()
	defer al.mu.Unlock()

	key := userID + ":" + action
	if lastTime, exists := al.rateLimit[key]; exists {
		if time.Since(lastTime) < time.Minute {
			return false
		}
	}

	activity := Activity{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now().UTC(),
		IPAddress: ip,
	}

	al.activities = append(al.activities, activity)
	al.rateLimit[key] = activity.Timestamp
	return true
}

func (al *ActivityLogger) GetActivities() []Activity {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return append([]Activity{}, al.activities...)
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		action := r.Method + " " + r.URL.Path
		ip := r.RemoteAddr

		al.LogActivity(userID, action, ip)

		next.ServeHTTP(w, r)
	})
}

func main() {
	logger := NewActivityLogger()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/admin/activities", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		activities := logger.GetActivities()
		json.NewEncoder(w).Encode(activities)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: logger.Middleware(mux),
	}

	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}