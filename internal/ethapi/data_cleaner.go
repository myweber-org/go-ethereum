package main

import "fmt"

func removeDuplicates(input []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    for _, value := range input {
        if !seen[value] {
            seen[value] = true
            result = append(result, value)
        }
    }
    return result
}

func main() {
    data := []int{1, 2, 2, 3, 4, 4, 5}
    cleaned := removeDuplicates(data)
    fmt.Println("Original:", data)
    fmt.Println("Cleaned:", cleaned)
}package utils

import (
	"regexp"
	"strings"
)

func SanitizeString(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	
	return cleaned
}