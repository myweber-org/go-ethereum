
package main

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
    input := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := RemoveDuplicates(input)
    fmt.Println("Original:", input)
    fmt.Println("Cleaned:", cleaned)
}package utils

import (
	"regexp"
	"strings"
)

func SanitizeString(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	
	return cleaned
}

func NormalizeWhitespace(input string) string {
	return SanitizeString(input)
}