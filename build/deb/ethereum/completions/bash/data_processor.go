package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Timestamp string `json:"timestamp"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func TransformUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func ProcessUserData(rawData []byte) (*UserProfile, error) {
	var profile UserProfile
	err := json.Unmarshal(rawData, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if profile.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", profile.ID)
	}

	profile.Username = TransformUsername(profile.Username)

	if !ValidateEmail(profile.Email) {
		return nil, fmt.Errorf("invalid email format: %s", profile.Email)
	}

	if profile.Age < 0 || profile.Age > 120 {
		return nil, fmt.Errorf("age out of valid range: %d", profile.Age)
	}

	return &profile, nil
}

func main() {
	jsonData := []byte(`{
		"id": 1001,
		"username": "  JohnDoe  ",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"timestamp": "2024-01-15T10:30:00Z"
	}`)

	processedProfile, err := ProcessUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed User Profile:\n")
	fmt.Printf("ID: %d\n", processedProfile.ID)
	fmt.Printf("Username: %s\n", processedProfile.Username)
	fmt.Printf("Email: %s\n", processedProfile.Email)
	fmt.Printf("Age: %d\n", processedProfile.Age)
	fmt.Printf("Active: %t\n", processedProfile.Active)
	fmt.Printf("Timestamp: %s\n", processedProfile.Timestamp)
}