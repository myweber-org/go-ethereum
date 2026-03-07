package data

import (
	"regexp"
	"strings"
)

// CleanString removes extra whitespace and trims the input string
func CleanString(input string) string {
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	return strings.TrimSpace(cleaned)
}

// ValidateEmail checks if the provided string is a valid email format
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}
	return matched
}

// FilterSlice removes empty strings from a slice
func FilterSlice(slice []string) []string {
	var result []string
	for _, item := range slice {
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}