package data_processor

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	var builder strings.Builder
	for _, r := range input {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func TransformToSlug(text string) string {
	text = strings.ToLower(text)
	text = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(text, "-")
	text = strings.Trim(text, "-")
	return text
}

func ExtractDomain(email string) (string, error) {
	if !ValidateEmail(email) {
		return "", errors.New("invalid email format")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", errors.New("malformed email address")
	}
	return parts[1], nil
}