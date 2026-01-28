package csvutils

import (
	"strings"
	"unicode"
)

// SanitizeField removes leading/trailing whitespace and replaces
// problematic characters in CSV fields
func SanitizeField(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace newlines and carriage returns with spaces
	replaced := strings.ReplaceAll(trimmed, "\n", " ")
	replaced = strings.ReplaceAll(replaced, "\r", " ")
	
	// Remove any remaining control characters
	var result strings.Builder
	for _, r := range replaced {
		if unicode.IsGraphic(r) || r == ' ' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// NormalizeWhitespace collapses multiple whitespace characters into single spaces
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
}package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 2, 3, 4, 4, 5, 5, 1}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
}