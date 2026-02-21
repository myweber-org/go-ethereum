
package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateAndNormalize(data UserData) (UserData, error) {
	var normalized UserData

	normalized.Username = strings.TrimSpace(data.Username)
	if normalized.Username == "" {
		return UserData{}, fmt.Errorf("username cannot be empty")
	}

	normalized.Email = strings.ToLower(strings.TrimSpace(data.Email))
	if !strings.Contains(normalized.Email, "@") {
		return UserData{}, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}
	normalized.Age = data.Age

	return normalized, nil
}

func main() {
	testData := UserData{
		Username: "  JohnDoe  ",
		Email:    "  TEST@EXAMPLE.COM",
		Age:      25,
	}

	result, err := ValidateAndNormalize(testData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Normalized Data: %+v\n", result)
}