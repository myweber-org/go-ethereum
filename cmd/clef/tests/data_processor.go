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

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func TransformUserData(data UserData) (UserData, error) {
	if err := ValidateEmail(data.Email); err != nil {
		return UserData{}, err
	}

	transformed := UserData{
		Email:    strings.ToLower(data.Email),
		Username: SanitizeUsername(data.Username),
		Age:      data.Age,
	}

	if transformed.Age < 0 {
		transformed.Age = 0
	}

	return transformed, nil
}