
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
}
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

func (dc *DataCleaner) IsDuplicate(item string) bool {
	normalized := dc.Normalize(item)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) Clean(items []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, item := range items {
		if !dc.IsDuplicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"Apple", "apple ", " Banana", "banana", "Cherry"}
	cleaned := cleaner.Clean(data)
	fmt.Println("Cleaned data:", cleaned)
}