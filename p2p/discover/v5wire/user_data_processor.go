
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}

func ValidateUser(u User) error {
	if u.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", u.ID)
	}

	if len(strings.TrimSpace(u.Username)) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("invalid email format: %s", u.Email)
	}

	return nil
}

func ParseUserJSON(data []byte) (*User, error) {
	var user User
	err := json.Unmarshal(data, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user JSON: %w", err)
	}

	err = ValidateUser(user)
	if err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	return &user, nil
}

func FormatUserOutput(u User) string {
	status := "inactive"
	if u.Active {
		status = "active"
	}
	return fmt.Sprintf("User #%d: %s (%s) - %s", u.ID, u.Username, u.Email, status)
}