
package main

import (
	"regexp"
	"strings"
)

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}

func SanitizeInput(input string) string {
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}