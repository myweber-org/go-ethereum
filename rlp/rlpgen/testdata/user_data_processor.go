package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Password string
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}

func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasUpper && hasLower && hasDigit
}

func ProcessUserData(data UserData) (UserData, error) {
	if !ValidateUsername(data.Username) {
		return UserData{}, ErrInvalidUsername
	}
	data.Email = SanitizeEmail(data.Email)
	if !ValidatePassword(data.Password) {
		return UserData{}, ErrWeakPassword
	}
	return data, nil
}

var (
	ErrInvalidUsername = errors.New("invalid username format")
	ErrWeakPassword    = errors.New("password does not meet security requirements")
)