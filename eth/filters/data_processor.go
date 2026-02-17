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
	var processed UserData

	if strings.TrimSpace(data.Username) == "" {
		return processed, fmt.Errorf("username cannot be empty")
	}
	processed.Username = strings.ToLower(strings.TrimSpace(data.Username))

	if !strings.Contains(data.Email, "@") {
		return processed, fmt.Errorf("invalid email format")
	}
	processed.Email = strings.ToLower(strings.TrimSpace(data.Email))

	if data.Age < 0 || data.Age > 150 {
		return processed, fmt.Errorf("age must be between 0 and 150")
	}
	processed.Age = data.Age

	return processed, nil
}

func main() {
	sampleData := UserData{
		Username: "  TestUser  ",
		Email:    "EXAMPLE@DOMAIN.COM",
		Age:      25,
	}

	result, err := ValidateAndTransform(sampleData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Processed Data: %+v\n", result)
}package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateJSON(data []byte) (*User, error) {
	var user User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if user.ID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.Name == "" {
		return nil, fmt.Errorf("user name cannot be empty")
	}
	if user.Email == "" {
		return nil, fmt.Errorf("user email cannot be empty")
	}

	return &user, nil
}

func main() {
	jsonData := []byte(`{"id": 123, "name": "John Doe", "email": "john@example.com"}`)
	user, err := ValidateJSON(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Validated user: %+v\n", user)
}