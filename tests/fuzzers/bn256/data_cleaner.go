
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

func (dc *DataCleaner) Deduplicate(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.validateItem(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) validateItem(item string) bool {
	return len(item) > 0 && len(item) <= 100
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"apple",
		"Apple",
		"banana",
		"  banana  ",
		"",
		"orange",
		"orange",
		strings.Repeat("x", 150),
	}
	
	cleaned := cleaner.Deduplicate(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
	fmt.Printf("Count: %d -> %d\n", len(data), len(cleaned))
	
	cleaner.Reset()
}