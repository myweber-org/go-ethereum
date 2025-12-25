package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Bio      string
}

func SanitizeInput(input string) string {
	// Remove leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(trimmed, " ")
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

func ProcessUserData(data UserData) (UserData, error) {
	sanitizedData := UserData{
		Username: SanitizeInput(data.Username),
		Email:    SanitizeInput(data.Email),
		Bio:      SanitizeInput(data.Bio),
	}

	if !ValidateUsername(sanitizedData.Username) {
		return UserData{}, &ValidationError{Field: "username", Message: "invalid username format"}
	}

	if !ValidateEmail(sanitizedData.Email) {
		return UserData{}, &ValidationError{Field: "email", Message: "invalid email format"}
	}

	if len(sanitizedData.Bio) > 500 {
		sanitizedData.Bio = sanitizedData.Bio[:500]
	}

	return sanitizedData, nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}