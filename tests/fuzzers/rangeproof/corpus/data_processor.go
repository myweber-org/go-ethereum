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

func ValidateAndTransform(data UserData) (UserData, error) {
	if strings.TrimSpace(data.Username) == "" {
		return data, fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return data, fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return data, fmt.Errorf("age must be between 0 and 150")
	}

	transformed := UserData{
		Username: strings.ToLower(strings.TrimSpace(data.Username)),
		Email:    strings.ToLower(strings.TrimSpace(data.Email)),
		Age:      data.Age,
	}
	return transformed, nil
}

func main() {
	sampleData := UserData{
		Username: "  TestUser  ",
		Email:    "TEST@EXAMPLE.COM",
		Age:      25,
	}

	result, err := ValidateAndTransform(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original: %+v\n", sampleData)
	fmt.Printf("Processed: %+v\n", result)
}