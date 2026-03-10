package main

import (
    "encoding/json"
    "log"
    "os"
    "time"
)

type UserActivity struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Timestamp time.Time `json:"timestamp"`
    Details   string    `json:"details,omitempty"`
}

func LogActivity(userID, action, details string) error {
    activity := UserActivity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(activity); err != nil {
        return err
    }

    return nil
}

func main() {
    err := LogActivity("user123", "login", "Successful authentication")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    err = LogActivity("user456", "file_upload", "Uploaded profile picture")
    if err != nil {
        log.Printf("Failed to log activity: %v", err)
    }

    log.Println("Activity logging completed")
}