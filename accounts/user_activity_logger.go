
package main

import (
    "encoding/json"
    "fmt"
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

func logActivity(userID, action, details string) {
    activity := UserActivity{
        UserID:    userID,
        Action:    action,
        Timestamp: time.Now().UTC(),
        Details:   details,
    }

    data, err := json.MarshalIndent(activity, "", "  ")
    if err != nil {
        log.Printf("Failed to marshal activity: %v", err)
        return
    }

    fmt.Println(string(data))

    file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Failed to open log file: %v", err)
        return
    }
    defer file.Close()

    if _, err := file.Write(append(data, '\n')); err != nil {
        log.Printf("Failed to write to log file: %v", err)
    }
}

func main() {
    logActivity("user123", "login", "Successful authentication")
    logActivity("user456", "purchase", "Order ID: 78910")
    logActivity("user123", "logout", "Session ended")
}