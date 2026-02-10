package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	
	re := regexp.MustCompile(`\s+`)
	trimmed = re.ReplaceAllString(trimmed, " ")
	
	var result strings.Builder
	for _, r := range trimmed {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}