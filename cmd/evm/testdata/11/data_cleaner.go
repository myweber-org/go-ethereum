package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	re := regexp.MustCompile(`\s+`)
	normalized := re.ReplaceAllString(trimmed, " ")
	
	// Remove potentially dangerous characters for SQL injection prevention
	re = regexp.MustCompile(`[;'"\\]`)
	cleaned := re.ReplaceAllString(normalized, "")
	
	return cleaned
}

func ValidateEmail(email string) bool {
	// Simple email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}