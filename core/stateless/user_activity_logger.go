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

func (l *ActivityLogger) LogActivity(userID, eventType, details string) error {
    event := ActivityEvent{
        UserID:    userID,
        EventType: eventType,
        Timestamp: time.Now(),
        Details:   details,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    data = append(data, '\n')
    _, err = l.logFile.Write(data)
    return err
}

func (l *ActivityLogger) Close() error {
    return l.logFile.Close()
}

func main() {
    logger, err := NewActivityLogger("user_activity.log")
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    events := []struct {
        userID    string
        eventType string
        details   string
    }{
        {"user123", "LOGIN", "User logged in from IP 192.168.1.100"},
        {"user123", "VIEW_PAGE", "Viewed dashboard"},
        {"user456", "REGISTER", "New user registration"},
        {"user123", "LOGOUT", "Session ended"},
    }

    for _, e := range events {
        err := logger.LogActivity(e.userID, e.eventType, e.details)
        if err != nil {
            fmt.Printf("Failed to log activity: %v\n", err)
        }
    }

    fmt.Println("Activity logging completed")
}