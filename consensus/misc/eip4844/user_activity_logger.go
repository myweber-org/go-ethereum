package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type ActivityEvent struct {
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
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
        UserID:    userID,
        EventType: eventType,
        Timestamp: time.Now().UTC(),
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
        return
    }

    err = logger.LogActivity("user123", "search", "Searched for 'golang tutorials'")
    if err != nil {
        fmt.Printf("Failed to log activity: %v\n", err)
        return
    }

    fmt.Println("Activities logged successfully")
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
	al.handler.ServeHTTP(w, r)
	duration := time.Since(start)

	log.Printf("Activity: %s %s from %s took %v",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		duration,
	)
}