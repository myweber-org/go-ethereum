
package data_processor

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim spaces
	cleaned = strings.TrimSpace(cleaned)
	
	// Convert to lowercase for normalization
	cleaned = strings.ToLower(cleaned)
	
	return cleaned
}

func NormalizeEmail(email string) string {
	cleaned := CleanInput(email)
	
	// Remove dots before @ for Gmail-like normalization
	parts := strings.Split(cleaned, "@")
	if len(parts) == 2 {
		localPart := strings.ReplaceAll(parts[0], ".", "")
		return localPart + "@" + parts[1]
	}
	
	return cleaned
}

func ValidateInput(input string, minLength int) bool {
	cleaned := CleanInput(input)
	return len(cleaned) >= minLength
}