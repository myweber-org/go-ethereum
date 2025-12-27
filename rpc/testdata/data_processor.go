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

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func main() {
	user := UserData{
		Username: "  TestUser  ",
		Email:    "test@example.com",
		Age:      25,
	}

	if err := ValidateUserData(user); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	user = TransformUsername(user)
	fmt.Printf("Processed user: %+v\n", user)
}