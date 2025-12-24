package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type UserActivity struct {
	UserID    string
	Action    string
	Timestamp time.Time
	SessionID string
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

func (al *ActivityLogger) LogActivity(userID, action, sessionID string) {
	activity := UserActivity{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now(),
		SessionID: sessionID,
	}

	logEntry := fmt.Sprintf("%s | %s | %s | %s\n",
		activity.Timestamp.Format(time.RFC3339),
		activity.UserID,
		activity.Action,
		activity.SessionID)

	if _, err := al.logFile.WriteString(logEntry); err != nil {
		log.Printf("Failed to write activity log: %v", err)
	}
}

func (al *ActivityLogger) Close() {
	if al.logFile != nil {
		al.logFile.Close()
	}
}

func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

func main() {
	logger, err := NewActivityLogger("user_activities.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	sessionID := generateSessionID()

	logger.LogActivity("user123", "login", sessionID)
	logger.LogActivity("user123", "view_profile", sessionID)
	logger.LogActivity("user456", "login", generateSessionID())
	logger.LogActivity("user123", "logout", sessionID)

	fmt.Println("Activity logging completed. Check user_activities.log")
}