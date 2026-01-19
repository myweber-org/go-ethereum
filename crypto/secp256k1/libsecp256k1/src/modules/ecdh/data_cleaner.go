package datautils

import (
	"regexp"
	"strings"
	"unicode"
)

// SanitizeString removes potentially harmful characters and normalizes whitespace
func SanitizeString(input string) string {
	if input == "" {
		return input
	}

	// Remove null characters and control characters except common whitespace
	cleaned := strings.Map(func(r rune) rune {
		if r == 0 || (unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r') {
			return -1
		}
		return r
	}, input)

	// Normalize multiple spaces/tabs to single space
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")

	// Trim leading/trailing whitespace
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

// NormalizeWhitespace converts all whitespace variations to standard spaces
func NormalizeWhitespace(input string) string {
	whitespaceRegex := regexp.MustCompile(`[\s\p{Z}]+`)
	return whitespaceRegex.ReplaceAllString(input, " ")
}

// IsSafeForDatabase checks if string contains only safe characters
func IsSafeForDatabase(input string) bool {
	safeRegex := regexp.MustCompile(`^[\p{L}\p{N}\p{P}\p{Z}]+$`)
	return safeRegex.MatchString(input)
}