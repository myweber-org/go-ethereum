package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

func SanitizeInput(input string) string {
	// Remove leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	// Escape potentially dangerous HTML characters
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(trimmed)
}

func ValidateEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailPattern, email)
	return err == nil && matched
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	usernamePattern := `^[a-zA-Z0-9_-]+$`
	matched, err := regexp.MatchString(usernamePattern, username)
	return err == nil && matched
}

func ProcessUserData(data UserData) (UserData, error) {
	sanitizedData := UserData{
		Username: SanitizeInput(data.Username),
		Email:    SanitizeInput(data.Email),
		Comments: SanitizeInput(data.Comments),
	}

	if !ValidateUsername(sanitizedData.Username) {
		return UserData{}, ErrInvalidUsername
	}

	if !ValidateEmail(sanitizedData.Email) {
		return UserData{}, ErrInvalidEmail
	}

	return sanitizedData, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrInvalidEmail    = errors.New("invalid email format")
)