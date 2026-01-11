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

func (al *ActivityLogger) LogActivity(userID, action, details string) error {
	logEntry := ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	jsonData = append(jsonData, '\n')
	_, err = al.logFile.Write(jsonData)
	return err
}

func (al *ActivityLogger) Close() error {
	return al.logFile.Close()
}

func main() {
	logger, err := NewActivityLogger("user_activity.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	err = logger.LogActivity("user123", "login", "Successful authentication")
	if err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	err = logger.LogActivity("user456", "file_upload", "Uploaded profile picture")
	if err != nil {
		log.Printf("Failed to log activity: %v", err)
	}

	fmt.Println("Activity logging completed")
}