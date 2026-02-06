package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

func ValidateAndTransform(data UserData) (UserData, error) {
	var processed UserData

	if data.Age < 0 || data.Age > 150 {
		return processed, errors.New("invalid age range")
	}
	processed.Age = data.Age

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(data.Email) {
		return processed, errors.New("invalid email format")
	}
	processed.Email = strings.ToLower(strings.TrimSpace(data.Email))

	username := strings.TrimSpace(data.Username)
	if len(username) < 3 || len(username) > 20 {
		return processed, errors.New("username must be 3-20 characters")
	}
	processed.Username = strings.ToLower(username)

	return processed, nil
}