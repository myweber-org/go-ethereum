package utils

import (
	"regexp"
	"strings"
)

// SanitizeInput cleans user-provided strings by removing excessive whitespace
// and trimming leading/trailing spaces. It replaces multiple spaces/tabs/newlines
// with a single space.
func SanitizeInput(input string) string {
	// Compile regex to match any whitespace sequence (spaces, tabs, newlines)
	re := regexp.MustCompile(`\s+`)
	// Replace all whitespace sequences with a single space
	cleaned := re.ReplaceAllString(input, " ")
	// Trim any remaining leading/trailing spaces
	return strings.TrimSpace(cleaned)
}