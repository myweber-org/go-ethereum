package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func ValidateUsername(username string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,20}$", username)
	return matched
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func TransformData(data UserData) (UserData, error) {
	if !ValidateUsername(data.Username) {
		return UserData{}, fmt.Errorf("invalid username format")
	}

	if !ValidateEmail(data.Email) {
		return UserData{}, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}

	transformed := data
	transformed.Username = strings.ToLower(data.Username)
	transformed.Email = strings.ToLower(data.Email)

	return transformed, nil
}

func ProcessJSONInput(jsonData []byte) (UserData, error) {
	var userData UserData
	err := json.Unmarshal(jsonData, &userData)
	if err != nil {
		return UserData{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return TransformData(userData)
}

func main() {
	jsonInput := `{"id": 1, "username": "TestUser_123", "email": "TEST@EXAMPLE.COM", "age": 25}`
	
	processedData, err := ProcessJSONInput([]byte(jsonInput))
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	fmt.Printf("Processed data: %+v\n", processedData)
}