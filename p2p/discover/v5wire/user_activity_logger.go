package middleware

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ActivityLogger struct {
	limiter *rate.Limiter
	logFunc func(userID string, action string, timestamp time.Time)
}

func NewActivityLogger(rps int, logFunc func(string, string, time.Time)) *ActivityLogger {
	return &ActivityLogger{
		limiter: rate.NewLimiter(rate.Limit(rps), rps),
		logFunc: logFunc,
	}
}

func (al *ActivityLogger) LogActivity(userID, action string) {
	if al.logFunc == nil {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	if err := al.limiter.Wait(ctx); err == nil {
		al.logFunc(userID, action, time.Now().UTC())
	}
}

func (al *ActivityLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := extractUserID(r)
		action := r.Method + " " + r.URL.Path
		
		go al.LogActivity(userID, action)
		
		next.ServeHTTP(w, r)
	})
}

func extractUserID(r *http.Request) string {
	if user := r.Header.Get("X-User-ID"); user != "" {
		return user
	}
	return "anonymous"
}