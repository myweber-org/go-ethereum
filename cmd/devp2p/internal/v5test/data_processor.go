
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
	return strings.ToLower(strings.TrimSpace(username))
}

func TransformProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	profile.Username = NormalizeUsername(profile.Username)

	if profile.Age < 0 || profile.Age > 150 {
		return profile, errors.New("age out of valid range")
	}

	return profile, nil
}

func ProcessProfiles(profiles []UserProfile) ([]UserProfile, []error) {
	var processed []UserProfile
	var errs []error

	for _, profile := range profiles {
		transformed, err := TransformProfile(profile)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, errs
}