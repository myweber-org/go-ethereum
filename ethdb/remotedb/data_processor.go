package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Email     string
	Username  string
	CreatedAt time.Time
	Active    bool
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserProfile(profile UserProfile) error {
	if profile.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if !emailRegex.MatchString(profile.Email) {
		return errors.New("invalid email format")
	}

	if len(strings.TrimSpace(profile.Username)) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if profile.CreatedAt.After(time.Now()) {
		return errors.New("creation date cannot be in the future")
	}

	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserData(profile UserProfile) (UserProfile, error) {
	if err := ValidateUserProfile(profile); err != nil {
		return UserProfile{}, err
	}

	profile.Username = TransformUsername(profile.Username)
	profile.Active = true

	return profile, nil
}