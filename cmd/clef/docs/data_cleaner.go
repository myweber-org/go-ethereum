
package utils

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned := spaceRegex.ReplaceAllString(trimmed, " ")
	
	return cleaned
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}