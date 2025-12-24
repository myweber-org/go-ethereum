
package main

import (
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Username string
	Email    string
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validUsername.MatchString(username)
}

func SanitizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	return email
}

func IsValidEmail(email string) bool {
	email = SanitizeEmail(email)
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ProcessUserInput(username, email string) (User, error) {
	if !ValidateUsername(username) {
		return User{}, ErrInvalidUsername
	}
	if !IsValidEmail(email) {
		return User{}, ErrInvalidEmail
	}
	cleanEmail := SanitizeEmail(email)
	return User{
		Username: username,
		Email:    cleanEmail,
	}, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrInvalidEmail    = errors.New("invalid email format")
)