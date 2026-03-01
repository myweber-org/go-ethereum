package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
}

func NewActivityLog(userID, action, resource string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}
}

func (al *ActivityLog) Save() error {
	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(al)
	if err != nil {
		return err
	}

	_, err = file.Write(append(data, '\n'))
	return err
}

func main() {
	logger := NewActivityLog("user123", "CREATE", "/api/document")
	if err := logger.Save(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Activity logged successfully")
}
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
		
		activity := ActivityLog{
			UserID:    extractUserID(r),
			IPAddress: r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
			Timestamp: start,
		}
		
		logActivity(activity)
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		log.Printf("Request completed in %v", duration)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return "authenticated_user"
	}
	return "anonymous"
}

func logActivity(activity ActivityLog) {
	log.Printf("Activity: User=%s IP=%s %s %s at %s",
		activity.UserID,
		activity.IPAddress,
		activity.Method,
		activity.Path,
		activity.Timestamp.Format(time.RFC3339),
	)
}