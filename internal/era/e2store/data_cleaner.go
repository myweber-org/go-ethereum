package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func CleanInput(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized := spaceRegex.ReplaceAllString(trimmed, " ")
	
	// Remove non-printable characters
	var result strings.Builder
	for _, r := range normalized {
		if unicode.IsPrint(r) {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

func NormalizeWhitespace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}