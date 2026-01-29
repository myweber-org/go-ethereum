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
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := normalizeSpaces(trimmed)
	sanitized := removeSpecialChars(normalized)
	return strings.ToValidUTF8(sanitized, "")
}

func normalizeSpaces(s string) string {
	spaceRegex := regexp.MustCompile(`\s+`)
	return spaceRegex.ReplaceAllString(s, " ")
}

func removeSpecialChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || r == '-' || r == '_' || r == '.' {
			return r
		}
		return -1
	}, s)
}