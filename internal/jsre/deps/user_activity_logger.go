package main

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
	logger, err := NewActivityLogger("activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	events := []struct {
		userID    string
		eventType string
		details   string
	}{
		{"user123", "login", "Successful login from IP 192.168.1.100"},
		{"user456", "purchase", "Purchased item SKU-789 with total $49.99"},
		{"user123", "logout", "User logged out after 30 minutes"},
		{"user789", "search", "Searched for 'golang tutorials'"},
	}

	for _, e := range events {
		err := logger.LogActivity(e.userID, e.eventType, e.details)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		} else {
			fmt.Printf("Logged event: %s for user %s\n", e.eventType, e.userID)
		}
	}

	fmt.Println("Activity logging completed")
}