package datautils

func RemoveDuplicates[T comparable](slice []T) []T {
    seen := make(map[T]bool)
    result := []T{}
    
    for _, item := range slice {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    
    return result
}
package main

import (
	"fmt"
	"strings"
)

func deduplicateStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func normalizeWhitespace(input string) string {
	words := strings.Fields(input)
	return strings.Join(words, " ")
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana"}
	uniqueData := deduplicateStrings(data)
	fmt.Println("Deduplicated:", uniqueData)

	text := "  This   is    a   test   string  with  extra   spaces.  "
	cleanedText := normalizeWhitespace(text)
	fmt.Println("Normalized:", cleanedText)
}