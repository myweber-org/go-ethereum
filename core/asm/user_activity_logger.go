package middleware

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	limiter  *rate.Limiter
	logStore ActivityStore
}

type ActivityStore interface {
	LogActivity(ctx context.Context, userID string, action string, metadata map[string]interface{}) error
}

func NewActivityLogger(store ActivityStore, requestsPerSecond int) *ActivityLogger {
	return &ActivityLogger{
		limiter:  rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond*2),
		logStore: store,
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		if !al.limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		userID := extractUserID(r)
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		action := r.Method + " " + r.URL.Path
		metadata := map[string]interface{}{
			"user_agent": r.UserAgent(),
			"ip_address": r.RemoteAddr,
			"timestamp":  time.Now().UTC(),
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			
			if err := al.logStore.LogActivity(ctx, userID, action, metadata); err != nil {
				logError("failed to log activity", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		return parseToken(authHeader)
	}
	return ""
}

func parseToken(token string) string {
	return token
}

func logError(msg string, err error) {
}