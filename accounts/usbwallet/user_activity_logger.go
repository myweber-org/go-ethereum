package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ActivityType string

const (
	Login    ActivityType = "LOGIN"
	Logout   ActivityType = "LOGOUT"
	Purchase ActivityType = "PURCHASE"
	View     ActivityType = "VIEW"
)

type UserActivity struct {
	UserID    string       `json:"user_id"`
	Action    ActivityType `json:"action"`
	Timestamp time.Time    `json:"timestamp"`
	Details   string       `json:"details,omitempty"`
}

func NewUserActivity(userID string, action ActivityType, details string) *UserActivity {
	return &UserActivity{
		UserID:    userID,
		Action:    action,
		Timestamp: time.Now().UTC(),
		Details:   details,
	}
}

func (ua *UserActivity) Log() error {
	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(ua)
}

func main() {
	activities := []*UserActivity{
		NewUserActivity("user123", Login, "Successful login from Chrome"),
		NewUserActivity("user123", View, "Viewed product catalog"),
		NewUserActivity("user123", Purchase, "Bought item: laptop"),
		NewUserActivity("user456", Login, "Mobile app login"),
	}

	for _, activity := range activities {
		if err := activity.Log(); err != nil {
			log.Printf("Failed to log activity: %v", err)
		} else {
			fmt.Printf("Logged: %s - %s\n", activity.UserID, activity.Action)
		}
	}
}