package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func TransformData(rawJSON string) (*UserData, error) {
	var data UserData
	err := json.Unmarshal([]byte(rawJSON), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if !ValidateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = SanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawData := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	processed, err := TransformData(rawData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Processed Data: %+v\n", processed)
}