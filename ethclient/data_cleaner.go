
package main

import (
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range input {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}

func main() {
	data := []string{" apple", "banana ", "  apple", "banana", "", "  cherry  "}
	cleaned := CleanData(data)
	fmt.Println("Cleaned Data:", cleaned)
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

func (dc *DataCleaner) AddItem(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return false
	}
	dc.seen[normalized] = true
	return true
}

func (dc *DataCleaner) UniqueCount() int {
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
	
	fmt.Printf("Total unique items: %d\n", cleaner.UniqueCount())
}