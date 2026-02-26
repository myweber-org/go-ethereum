package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type ActivityKey string

const UserActivityKey ActivityKey = "user_activity"

type UserActivity struct {
	UserID    string
	Action    string
	Timestamp time.Time
	IPAddress string
	UserAgent string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		activity := UserActivity{
			UserID:    extractUserID(r),
			Action:    r.Method + " " + r.URL.Path,
			Timestamp: time.Now().UTC(),
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		ctx := context.WithValue(r.Context(), UserActivityKey, activity)
		next.ServeHTTP(w, r.WithContext(ctx))

		logActivity(activity)
	})
}

func extractUserID(r *http.Request) string {
	if auth := r.Header.Get("Authorization"); auth != "" {
		return "authenticated_user"
	}
	return "anonymous"
}

func logActivity(activity UserActivity) {
	log.Printf("Activity: User=%s Action=%s IP=%s Time=%s",
		activity.UserID,
		activity.Action,
		activity.IPAddress,
		activity.Timestamp.Format(time.RFC3339))
}

func GetActivityFromContext(ctx context.Context) (UserActivity, bool) {
	activity, ok := ctx.Value(UserActivityKey).(UserActivity)
	return activity, ok
}
package main

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Action    string
	Path      string
	Timestamp time.Time
}

var activityLogs []ActivityLog

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}

		logEntry := ActivityLog{
			UserID:    userID,
			Action:    r.Method,
			Path:      r.URL.Path,
			Timestamp: time.Now(),
		}

		activityLogs = append(activityLogs, logEntry)
		log.Printf("Activity: %s %s by %s", logEntry.Action, logEntry.Path, logEntry.UserID)

		next.ServeHTTP(w, r)
	})
}

func getActivityLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	for _, entry := range activityLogs {
		logLine := entry.Timestamp.Format("2006-01-02 15:04:05") + " - " + entry.UserID + " - " + entry.Action + " " + entry.Path + "\n"
		w.Write([]byte(logLine))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin/logs", getActivityLogs)

	loggedMux := loggingMiddleware(mux)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", loggedMux); err != nil {
		log.Fatal(err)
	}
}