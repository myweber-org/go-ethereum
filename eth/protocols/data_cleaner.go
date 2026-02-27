
package main

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := normalizeSpaces(trimmed)
	cleaned := removeSpecialChars(normalized)
	return cleaned
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