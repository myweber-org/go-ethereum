
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) AddItem(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return false
	}
	dc.seen[normalized] = true
	return true
}

func (dc *DataCleaner) GetUniqueCount() int {
	return len(dc.seen)
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	samples := []string{"  Apple  ", "apple", "BANANA", "banana ", "Cherry"}
	
	for _, item := range samples {
		if cleaner.AddItem(item) {
			fmt.Printf("Added: '%s'\n", item)
		} else {
			fmt.Printf("Duplicate skipped: '%s'\n", item)
		}
	}
	
	fmt.Printf("Unique items: %d\n", cleaner.GetUniqueCount())
}package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeString(input string) string {
	// Trim whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove extra internal whitespace
	space := regexp.MustCompile(`\s+`)
	normalized := space.ReplaceAllString(trimmed, " ")
	
	// Remove non-printable characters
	var result strings.Builder
	for _, r := range normalized {
		if unicode.IsPrint(r) {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

func NormalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), " ")
}package main

import "fmt"

func RemoveDuplicates(input []int) []int {
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
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}