
package main

import (
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Username string
	Email    string
	Age      int
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}

func SanitizeEmail(email string) string {
	trimmed := strings.TrimSpace(email)
	return strings.ToLower(trimmed)
}

func ValidateUserAge(age int) bool {
	return age >= 18 && age <= 120
}

func ProcessUserInput(username, email string, age int) (User, error) {
	if !ValidateUsername(username) {
		return User{}, ErrInvalidUsername
	}

	if !ValidateUserAge(age) {
		return User{}, ErrInvalidAge
	}

	sanitizedEmail := SanitizeEmail(email)

	user := User{
		Username: username,
		Email:    sanitizedEmail,
		Age:      age,
	}

	return user, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrInvalidAge      = errors.New("age must be between 18 and 120")
)