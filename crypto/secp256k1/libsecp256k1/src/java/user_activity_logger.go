package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ActivityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s %d %v",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
		)
	})
}package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityEvent struct {
    Timestamp time.Time `json:"timestamp"`
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Details   string    `json:"details"`
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

func (l *ActivityLogger) LogActivity(userID, eventType, details string) error {
    event := ActivityEvent{
        Timestamp: time.Now(),
        UserID:    userID,
        EventType: eventType,
        Details:   details,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    _, err = l.logFile.Write(append(data, '\n'))
    return err
}

func (l *ActivityLogger) Close() error {
    return l.logFile.Close()
}

func main() {
    logger, err := NewActivityLogger("activity.log")
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    err = logger.LogActivity("user123", "login", "User logged in from web browser")
    if err != nil {
        fmt.Printf("Failed to log activity: %v\n", err)
    }

    err = logger.LogActivity("user123", "search", "Searched for 'golang tutorials'")
    if err != nil {
        fmt.Printf("Failed to log activity: %v\n", err)
    }

    fmt.Println("Activity logging completed")
}