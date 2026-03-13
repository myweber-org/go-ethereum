package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateUserData(data UserData) error {
	if data.Email == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Username == "" {
		return errors.New("username is required")
	}
	if len(data.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: TransformUsername(username),
		Age:      age,
	}

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}
package main

import (
	"encoding/json"
	"errors"
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

func ValidateProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return errors.New("invalid user ID")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	if !usernameRegex.MatchString(profile.Username) {
		return errors.New("username must be 3-20 alphanumeric characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func TransformProfile(profile UserProfile) UserProfile {
	transformed := profile
	transformed.Username = strings.ToLower(profile.Username)
	transformed.Email = strings.ToLower(profile.Email)
	transformed.Timestamp = strings.ReplaceAll(profile.Timestamp, " ", "T")
	return transformed
}

func ProcessUserData(inputJSON string) (string, error) {
	var profile UserProfile
	err := json.Unmarshal([]byte(inputJSON), &profile)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %v", err)
	}

	err = ValidateProfile(profile)
	if err != nil {
		return "", fmt.Errorf("validation failed: %v", err)
	}

	transformedProfile := TransformProfile(profile)
	outputJSON, err := json.MarshalIndent(transformedProfile, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return string(outputJSON), nil
}

func main() {
	sampleInput := `{
		"id": 1001,
		"username": "John_Doe",
		"email": "John@Example.COM",
		"age": 30,
		"active": true,
		"timestamp": "2024-01-15 14:30:00"
	}`

	result, err := ProcessUserData(sampleInput)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Processed profile:")
	fmt.Println(result)
}