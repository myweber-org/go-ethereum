package middleware

import (
	"log"
	"net/http"
	"time"
)

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	writer := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	
	al.handler.ServeHTTP(writer, r)
	
	duration := time.Since(start)
	log.Printf("%s %s %d %v", r.Method, r.URL.Path, writer.statusCode, duration)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type ActivityLogger struct {
	mu          sync.RWMutex
	userLimits  map[string]time.Time
	rateLimit   time.Duration
	logFilePath string
}

func NewActivityLogger(limit time.Duration, logFile string) *ActivityLogger {
	return &ActivityLogger{
		userLimits:  make(map[string]time.Time),
		rateLimit:   limit,
		logFilePath: logFile,
	}
}

func (al *ActivityLogger) LogActivity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIP := r.RemoteAddr
		path := r.URL.Path
		method := r.Method
		timestamp := time.Now().UTC()

		al.mu.RLock()
		lastTime, exists := al.userLimits[userIP]
		al.mu.RUnlock()

		if exists && timestamp.Sub(lastTime) < al.rateLimit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		al.mu.Lock()
		al.userLimits[userIP] = timestamp
		al.mu.Unlock()

		logEntry := timestamp.Format("2006-01-02 15:04:05") + " | " +
			userIP + " | " + method + " " + path

		go func(entry string) {
			log.Println(entry)
		}(logEntry)

		next.ServeHTTP(w, r)
	})
}