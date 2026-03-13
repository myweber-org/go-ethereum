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

    data = append(data, '\n')
    _, err = l.logFile.Write(data)
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

    events := []struct {
        userID, eventType, details string
    }{
        {"user123", "login", "User logged in from Chrome browser"},
        {"user123", "view_page", "Viewed dashboard page"},
        {"user456", "purchase", "Purchased item SKU-789"},
        {"user123", "logout", "User logged out after 15 minutes"},
    }

    for _, e := range events {
        if err := logger.LogActivity(e.userID, e.eventType, e.details); err != nil {
            fmt.Printf("Failed to log activity: %v\n", err)
        }
    }

    fmt.Println("Activity logging completed successfully")
}