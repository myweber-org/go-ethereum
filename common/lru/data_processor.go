package main

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

type UserProfile struct {
	Username string
	Email    string
	Age      int
}

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validPattern.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func NormalizeEmail(email string) (string, error) {
	trimmed := strings.TrimSpace(email)
	lower := strings.ToLower(trimmed)
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(lower) {
		return "", errors.New("invalid email format")
	}
	return lower, nil
}

func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateUsername(profile.Username); err != nil {
		return UserProfile{}, err
	}
	normalizedEmail, err := NormalizeEmail(profile.Email)
	if err != nil {
		return UserProfile{}, err
	}
	if err := ValidateAge(profile.Age); err != nil {
		return UserProfile{}, err
	}
	return UserProfile{
		Username: profile.Username,
		Email:    normalizedEmail,
		Age:      profile.Age,
	}, nil
}

func SanitizeString(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, input)
}