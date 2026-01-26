package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type activityKey string

const UserActivityKey activityKey = "user_activity"

type UserActivity struct {
	UserID    string
	Action    string
	Timestamp time.Time
	IPAddress string
	UserAgent string
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		activity := UserActivity{
			Timestamp: start,
			IPAddress: r.RemoteAddr,
			UserAgent: r.UserAgent(),
			Action:    r.Method + " " + r.URL.Path,
		}

		ctx := context.WithValue(r.Context(), UserActivityKey, &activity)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("Activity: %+v, Duration: %v", activity, duration)
	})
}

func SetActivityUser(ctx context.Context, userID string) {
	if activity, ok := ctx.Value(UserActivityKey).(*UserActivity); ok {
		activity.UserID = userID
	}
}