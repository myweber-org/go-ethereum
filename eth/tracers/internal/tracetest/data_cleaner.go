package csvutils

import (
	"strings"
	"unicode"
)

// SanitizeField removes leading/trailing whitespace and replaces
// problematic characters in CSV fields
func SanitizeField(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace newlines and carriage returns with spaces
	replaced := strings.ReplaceAll(trimmed, "\n", " ")
	replaced = strings.ReplaceAll(replaced, "\r", " ")
	
	// Remove any remaining control characters
	var result strings.Builder
	for _, r := range replaced {
		if unicode.IsGraphic(r) || r == ' ' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// NormalizeWhitespace collapses multiple whitespace characters into single spaces
func NormalizeWhitespace(input string) string {
	var result strings.Builder
	prevSpace := false
	
	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}
	
	return result.String()
}