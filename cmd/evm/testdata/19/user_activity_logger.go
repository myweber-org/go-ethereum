
package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	Timestamp  time.Time
	Method     string
	Path       string
	RemoteAddr string
	UserAgent  string
	StatusCode int
	Duration   time.Duration
}

type ActivityLogger struct {
	logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{logger: logger}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(recorder, r)
		
		duration := time.Since(start)
		activity := ActivityLog{
			Timestamp:  start,
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			StatusCode: recorder.statusCode,
			Duration:   duration,
		}
		
		al.logger.Printf(
			"%s %s %s %d %s %s",
			activity.Timestamp.Format(time.RFC3339),
			activity.Method,
			activity.Path,
			activity.StatusCode,
			activity.RemoteAddr,
			activity.Duration,
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}