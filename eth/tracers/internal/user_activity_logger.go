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
	
	log.Printf("%s %s %d %v",
		r.Method,
		r.URL.Path,
		recorder.statusCode,
		duration,
	)
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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityEvent struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details,omitempty"`
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

func (al *ActivityLogger) LogActivity(userID, eventType, details string) error {
	event := ActivityEvent{
		UserID:    userID,
		EventType: eventType,
		Timestamp: time.Now(),
		Details:   details,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	eventJSON = append(eventJSON, '\n')
	_, err = al.logFile.Write(eventJSON)
	return err
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("user_activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	err = logger.LogActivity("user123", "login", "User logged in from web browser")
	if err != nil {
		log.Fatal(err)
	}

	err = logger.LogActivity("user123", "view_page", "Viewed dashboard")
	if err != nil {
		log.Fatal(err)
	}

	err = logger.LogActivity("user456", "logout", "Session expired")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Activity logged successfully")
}