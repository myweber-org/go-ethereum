package csvutil

import (
	"strings"
)

// CleanCSVRow removes leading/trailing whitespace from each field
// and filters out completely empty rows
func CleanCSVRow(row []string) []string {
	cleaned := make([]string, 0, len(row))
	allEmpty := true

	for _, field := range row {
		trimmed := strings.TrimSpace(field)
		cleaned = append(cleaned, trimmed)
		if trimmed != "" {
			allEmpty = false
		}
	}

	if allEmpty {
		return []string{}
	}
	return cleaned
}

// CleanCSVData processes multiple rows and returns only non-empty rows
func CleanCSVData(data [][]string) [][]string {
	result := make([][]string, 0, len(data))
	for _, row := range data {
		cleaned := CleanCSVRow(row)
		if len(cleaned) > 0 {
			result = append(result, cleaned)
		}
	}
	return result
}
package main

import (
	"fmt"
	"strings"
)

func DeduplicateStrings(slice []string) []string {
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

func NormalizeWhitespace(input string) string {
	words := strings.Fields(input)
	return strings.Join(words, " ")
}

func main() {
	data := []string{"apple", "banana", "apple", "cherry", "banana"}
	unique := DeduplicateStrings(data)
	fmt.Println("Deduplicated:", unique)

	text := "   This   is    a   test   string   with   extra   spaces.   "
	normalized := NormalizeWhitespace(text)
	fmt.Println("Normalized:", normalized)
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

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) Deduplicate(values []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, v := range values {
		if !dc.IsDuplicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", "Apple ", "banana", "BANANA", " cherry", "cherry "}
	
	fmt.Println("Original data:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated data:", deduped)
	
	testValue := "  APPLE  "
	fmt.Printf("Is '%s' duplicate? %v\n", testValue, cleaner.IsDuplicate(testValue))
}