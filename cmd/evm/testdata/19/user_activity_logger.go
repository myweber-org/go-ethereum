
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
}package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLog struct {
	UserID    string
	Endpoint  string
	Method    string
	Timestamp time.Time
	IPAddress string
}

var activityChannel = make(chan ActivityLog, 100)

func init() {
	go processActivityLogs()
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		userID := "anonymous"
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			userID = extractUserID(authHeader)
		}

		activity := ActivityLog{
			UserID:    userID,
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Timestamp: start,
			IPAddress: r.RemoteAddr,
		}

		select {
		case activityChannel <- activity:
		default:
			log.Println("Activity log channel full, dropping log entry")
		}

		next.ServeHTTP(w, r)
	})
}

func extractUserID(token string) string {
	return "user_" + token[:8]
}

func processActivityLogs() {
	for activity := range activityChannel {
		log.Printf("ACTIVITY: %s %s %s %s %v\n",
			activity.UserID,
			activity.IPAddress,
			activity.Method,
			activity.Endpoint,
			activity.Timestamp.Format(time.RFC3339))
	}
}