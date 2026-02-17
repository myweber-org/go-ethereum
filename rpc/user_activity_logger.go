package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityEvent struct {
	Timestamp time.Time
	UserID    string
	EventType string
	Details   map[string]string
}

type ActivityLogger struct {
	logFile *os.File
}

func NewActivityLogger(filename string) (*ActivityLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &ActivityLogger{logFile: file}, nil
}

func (al *ActivityLogger) LogActivity(userID, eventType string, details map[string]string) {
	event := ActivityEvent{
		Timestamp: time.Now(),
		UserID:    userID,
		EventType: eventType,
		Details:   details,
	}

	logEntry := fmt.Sprintf("[%s] User: %s | Event: %s | Details: %v\n",
		event.Timestamp.Format("2006-01-02 15:04:05"),
		event.UserID,
		event.EventType,
		event.Details)

	if _, err := al.logFile.WriteString(logEntry); err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
}

func (al *ActivityLogger) Close() {
	if al.logFile != nil {
		al.logFile.Close()
	}
}

func main() {
	logger, err := NewActivityLogger("user_activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	details := map[string]string{
		"page":     "/dashboard",
		"action":   "button_click",
		"element":  "refresh_button",
	}

	logger.LogActivity("user123", "page_interaction", details)

	loginDetails := map[string]string{
		"ip_address": "192.168.1.100",
		"user_agent": "Mozilla/5.0",
	}

	logger.LogActivity("user456", "login", loginDetails)

	fmt.Println("Activity logging completed")
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
		"Method: %s | Path: %s | Status: %d | Duration: %v | RemoteAddr: %s | UserAgent: %s",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
		r.RemoteAddr,
		r.UserAgent(),
	)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}