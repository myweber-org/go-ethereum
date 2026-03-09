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
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	al.handler.ServeHTTP(recorder, r)

	duration := time.Since(start)
	log.Printf(
		"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
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

type ActivityLogger struct {
	handler http.Handler
}

func NewActivityLogger(handler http.Handler) *ActivityLogger {
	return &ActivityLogger{handler: handler}
}

func (al *ActivityLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
	
	al.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)
	
	log.Printf(
		"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
	)
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

type ActivityLogger struct {
	Logger *log.Logger
}

func NewActivityLogger(logger *log.Logger) *ActivityLogger {
	return &ActivityLogger{Logger: logger}
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
		
		al.Logger.Printf(
			"%s %s %d %s %s",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration.String(),
			r.RemoteAddr,
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
}package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time
	SessionID string
	UserID    string
	Action    string
	Details   string
}

func NewActivityLog(sessionID, userID, action, details string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now(),
		SessionID: sessionID,
		UserID:    userID,
		Action:    action,
		Details:   details,
	}
}

func (al *ActivityLog) String() string {
	return fmt.Sprintf("[%s] Session:%s User:%s Action:%s Details:%s",
		al.Timestamp.Format("2006-01-02 15:04:05"),
		al.SessionID,
		al.UserID,
		al.Action,
		al.Details)
}

type ActivityLogger struct {
	logFile *os.File
	logger  *log.Logger
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &ActivityLogger{
		logFile: file,
		logger:  log.New(file, "", 0),
	}, nil
}

func (al *ActivityLogger) LogActivity(logEntry *ActivityLog) {
	al.logger.Println(logEntry.String())
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("user_activity.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	sessionID := "sess_12345"
	userID := "user_67890"

	activities := []*ActivityLog{
		NewActivityLog(sessionID, userID, "LOGIN", "User logged in from IP 192.168.1.100"),
		NewActivityLog(sessionID, userID, "VIEW_PAGE", "Accessed dashboard page"),
		NewActivityLog(sessionID, userID, "UPDATE_SETTINGS", "Changed notification preferences"),
		NewActivityLog(sessionID, userID, "LOGOUT", "User logged out successfully"),
	}

	for _, activity := range activities {
		logger.LogActivity(activity)
		fmt.Println(activity.String())
	}

	fmt.Println("Activity logging completed.")
}