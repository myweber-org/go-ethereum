
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package datautils

import (
	"strings"
	"unicode"
)

// SanitizeString cleans up CSV field values by:
// - Removing leading/trailing whitespace
// - Trimming matching quotes (both single and double)
// - Collapsing multiple internal spaces to single space
// - Removing non-printable characters except standard whitespace
func SanitizeString(input string) string {
	if input == "" {
		return input
	}

	// Remove surrounding quotes if they match
	input = strings.TrimSpace(input)
	if len(input) >= 2 {
		firstChar := input[0]
		lastChar := input[len(input)-1]
		if (firstChar == '"' && lastChar == '"') || (firstChar == '\'' && lastChar == '\'') {
			input = input[1 : len(input)-1]
		}
	}

	// Collapse multiple spaces and remove non-printable characters
	var result strings.Builder
	prevSpace := false
	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace && r == ' ' {
				result.WriteRune(' ')
				prevSpace = true
			}
			continue
		}
		if unicode.IsPrint(r) {
			result.WriteRune(r)
			prevSpace = false
		}
	}

	return strings.TrimSpace(result.String())
}