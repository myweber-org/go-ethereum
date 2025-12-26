package main

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
    Details   string    `json:"details"`
}

func NewActivityLog(userID, action, details string) *ActivityLog {
    return &ActivityLog{
        Timestamp: time.Now(),
        UserID:    userID,
        Action:    action,
        Details:   details,
    }
}

func (al *ActivityLog) SaveToFile(filename string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    data, err := json.Marshal(al)
    if err != nil {
        return err
    }

    _, err = file.Write(append(data, '\n'))
    return err
}

func main() {
    logger := NewActivityLog("user123", "login", "User logged in from web browser")
    
    if err := logger.SaveToFile("activity.log"); err != nil {
        log.Fatal("Failed to save activity log:", err)
    }
    
    fmt.Println("Activity logged successfully")
}