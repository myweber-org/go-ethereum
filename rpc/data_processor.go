
package data_processor

import (
	"regexp"
	"strings"
	"unicode"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := normalizeSpaces(trimmed)
	return removeSpecialChars(normalized)
}

func normalizeSpaces(s string) string {
	spaceRegex := regexp.MustCompile(`\s+`)
	return spaceRegex.ReplaceAllString(s, " ")
}

func removeSpecialChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)
}

func Tokenize(s string) []string {
	cleaned := CleanInput(s)
	if cleaned == "" {
		return []string{}
	}
	return strings.Split(cleaned, " ")
}