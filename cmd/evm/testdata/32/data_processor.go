
package main

import (
	"strings"
	"unicode"
)

func CleanString(input string) string {
	return strings.TrimSpace(input)
}

func NormalizeWhitespace(input string) string {
	var result strings.Builder
	prevSpace := false

	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}

	return result.String()
}

func ToLowercase(input string) string {
	return strings.ToLower(input)
}

func ProcessInput(input string) string {
	cleaned := CleanString(input)
	normalized := NormalizeWhitespace(cleaned)
	return ToLowercase(normalized)
}