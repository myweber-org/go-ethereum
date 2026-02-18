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
	Resource  string    `json:"resource"`
}

func NewActivityLog(userID, action, resource string) *ActivityLog {
	return &ActivityLog{
		Timestamp: time.Now().UTC(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
	}
}

func (al *ActivityLog) Save() error {
	file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	logger := NewActivityLog("user123", "CREATE", "/api/document")
	if err := logger.Save(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Activity logged successfully")
}