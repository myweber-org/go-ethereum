package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type EventType string

const (
	Login    EventType = "LOGIN"
	Logout   EventType = "LOGOUT"
	Purchase EventType = "PURCHASE"
	View     EventType = "VIEW"
)

type UserActivity struct {
	UserID    string
	Event     EventType
	Timestamp time.Time
	Metadata  map[string]string
}

func NewUserActivity(userID string, event EventType) *UserActivity {
	return &UserActivity{
		UserID:    userID,
		Event:     event,
		Timestamp: time.Now().UTC(),
		Metadata:  make(map[string]string),
	}
}

func (ua *UserActivity) AddMetadata(key, value string) {
	ua.Metadata[key] = value
}

func (ua *UserActivity) Log() {
	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	logEntry := fmt.Sprintf("[%s] User: %s, Event: %s",
		ua.Timestamp.Format(time.RFC3339),
		ua.UserID,
		ua.Event,
	)

	if len(ua.Metadata) > 0 {
		logEntry += fmt.Sprintf(", Metadata: %v", ua.Metadata)
	}

	logEntry += "\n"

	if _, err := file.WriteString(logEntry); err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
}

func main() {
	activity := NewUserActivity("user123", Login)
	activity.AddMetadata("ip", "192.168.1.100")
	activity.AddMetadata("user_agent", "Mozilla/5.0")
	activity.Log()

	purchaseActivity := NewUserActivity("user456", Purchase)
	purchaseActivity.AddMetadata("item_id", "prod_789")
	purchaseActivity.AddMetadata("amount", "49.99")
	purchaseActivity.Log()

	fmt.Println("Activity logging completed")
}