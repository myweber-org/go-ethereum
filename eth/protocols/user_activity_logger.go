package middleware

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
		startTime := time.Now()
		userAgent := r.UserAgent()
		clientIP := r.RemoteAddr
		method := r.Method
		path := r.URL.Path

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(startTime)
		status := recorder.statusCode

		al.Logger.Printf(
			"IP: %s | Method: %s | Path: %s | Status: %d | Duration: %v | Agent: %s",
			clientIP, method, path, status, duration, userAgent,
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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
}

func NewActivityLog(userID, action, resource, details string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Details:   details,
	}
}

func (al *ActivityLog) ToJSON() ([]byte, error) {
	return json.MarshalIndent(al, "", "  ")
}

func (al *ActivityLog) SaveToFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := al.ToJSON()
	if err != nil {
		return err
	}

	_, err = file.Write(append(jsonData, '\n'))
	return err
}

func main() {
	logger := NewActivityLog(
		"user_12345",
		"LOGIN",
		"authentication",
		"User logged in from IP 192.168.1.100",
	)

	jsonOutput, err := logger.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Activity Log:")
	fmt.Println(string(jsonOutput))

	err = logger.SaveToFile("activity_logs.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Log saved to activity_logs.json")
}package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityLog struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
}

func logActivity(userID, action, resource string) {
	activity := ActivityLog{
		Timestamp: time.Now(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}

	logData, err := json.MarshalIndent(activity, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal activity: %v", err)
		return
	}

	logFile, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer logFile.Close()

	if _, err := logFile.Write(append(logData, '\n')); err != nil {
		log.Printf("Failed to write to log file: %v", err)
	}
}

func main() {
	logActivity("user123", "CREATE", "/api/document")
	logActivity("user456", "READ", "/api/report")
	logActivity("user789", "UPDATE", "/api/profile")

	fmt.Println("Activity logging completed. Check activity.log file.")
}