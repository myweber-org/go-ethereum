
package main

import (
	"fmt"
	"strings"
)

// DataCleaner holds methods for cleaning string slices
type DataCleaner struct{}

// RemoveDuplicates removes duplicate entries from a slice of strings
func (dc DataCleaner) RemoveDuplicates(input []string) []string {
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

// TrimWhitespace removes leading and trailing spaces from each string
func (dc DataCleaner) TrimWhitespace(input []string) []string {
	result := make([]string, len(input))
	for i, item := range input {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func main() {
	cleaner := DataCleaner{}
	data := []string{"  apple ", "banana", "  apple ", " cherry", "banana "}

	fmt.Println("Original data:", data)

	trimmed := cleaner.TrimWhitespace(data)
	fmt.Println("After trimming:", trimmed)

	unique := cleaner.RemoveDuplicates(trimmed)
	fmt.Println("After deduplication:", unique)
}package main

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
}