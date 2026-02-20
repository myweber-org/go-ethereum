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

func NormalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", errors.New("invalid email format")
	}
	return email, nil
}

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	normalizedEmail, err := NormalizeEmail(profile.Email)
	if err != nil {
		return UserProfile{}, err
	}
	profile.Email = normalizedEmail

	if err := ValidateUsername(profile.Username); err != nil {
		return UserProfile{}, err
	}

	if profile.Age < 0 || profile.Age > 150 {
		return UserProfile{}, errors.New("age must be between 0 and 150")
	}

	return profile, nil
}