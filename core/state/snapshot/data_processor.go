
package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        string
	Email     string
	Username  string
	Age       int
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateProfile(profile UserProfile) error {
	if profile.ID == "" {
		return errors.New("ID cannot be empty")
	}

	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if len(profile.Username) < 3 || len(profile.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if profile.Age < 0 || profile.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func NormalizeProfile(profile UserProfile) UserProfile {
	normalized := profile
	normalized.Email = strings.ToLower(strings.TrimSpace(profile.Email))
	normalized.Username = strings.TrimSpace(profile.Username)
	return normalized
}

func ProcessUserData(profile UserProfile) (UserProfile, error) {
	normalized := NormalizeProfile(profile)

	if err := ValidateProfile(normalized); err != nil {
		return UserProfile{}, err
	}

	return normalized, nil
}