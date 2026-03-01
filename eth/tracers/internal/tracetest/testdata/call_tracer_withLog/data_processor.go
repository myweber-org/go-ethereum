
package main

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	var result strings.Builder
	for _, r := range username {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
		"&", "&amp;",
	)
	return replacer.Replace(input)
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("password must contain uppercase, lowercase and digit")
	}
	return nil
}