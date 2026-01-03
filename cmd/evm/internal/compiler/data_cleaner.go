package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput removes leading/trailing whitespace and special characters
func SanitizeInput(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove special characters except alphanumeric and basic punctuation
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s.,!?-]`)
	sanitized := reg.ReplaceAllString(trimmed, "")
	
	return sanitized
}