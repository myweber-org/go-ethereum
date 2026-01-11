package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Username  string
	Email     string
	BirthDate string
	CreatedAt time.Time
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

func ValidateEmail(email string) error {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func CalculateAge(birthDate string) (int, error) {
	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, errors.New("invalid date format, use YYYY-MM-DD")
	}
	now := time.Now()
	age := now.Year() - parsedDate.Year()
	if now.YearDay() < parsedDate.YearDay() {
		age--
	}
	if age < 0 {
		return 0, errors.New("birth date cannot be in the future")
	}
	return age, nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	transformedUsername := TransformUsername(profile.Username)
	profile.Username = transformedUsername

	if err := ValidateUsername(profile.Username); err != nil {
		return profile, err
	}

	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	_, err := CalculateAge(profile.BirthDate)
	if err != nil {
		return profile, err
	}

	return profile, nil
}