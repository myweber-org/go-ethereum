
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	Method    string
	Path      string
	IP        string
	UserAgent string
	Duration  time.Duration
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(lrw, r)
		
		duration := time.Since(start)
		
		activity := ActivityLog{
			Timestamp: time.Now(),
			Method:    r.Method,
			Path:      r.URL.Path,
			IP:        r.RemoteAddr,
			UserAgent: r.UserAgent(),
			Duration:  duration,
		}
		
		log.Printf("Activity: %s %s from %s - %d - %v",
			activity.Method,
			activity.Path,
			activity.IP,
			lrw.statusCode,
			activity.Duration,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}