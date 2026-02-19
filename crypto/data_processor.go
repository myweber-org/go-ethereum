
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

func NormalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	profile.Username = NormalizeUsername(profile.Username)
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = time.Now().UTC()
	}

	return profile, nil
}