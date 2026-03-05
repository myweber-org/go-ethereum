
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func SanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return UserProfile{}, err
	}

	profile.Username = SanitizeUsername(profile.Username)

	if err := ValidateAge(profile.Age); err != nil {
		return UserProfile{}, err
	}

	return profile, nil
}