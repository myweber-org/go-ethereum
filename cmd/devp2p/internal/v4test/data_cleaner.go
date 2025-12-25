package utils

import (
	"regexp"
	"strings"
)

func SanitizeString(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	
	return cleaned
}