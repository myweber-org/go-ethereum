package main

import (
	"regexp"
	"strings"
)

func NormalizeEmail(email string) (string, bool) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil || !matched {
		return "", false
	}
	return email, true
}

func ValidateUsername(username string) bool {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return false
	}
	return matched
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}