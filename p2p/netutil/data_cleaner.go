
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

func (dc *DataCleaner) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := strings.ToLower(trimmed)
	return normalized
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	cleaned := dc.CleanString(value)
	if dc.seen[cleaned] {
		return true
	}
	dc.seen[cleaned] = true
	return false
}

func (dc *DataCleaner) AddItem(value string) bool {
	cleaned := dc.CleanString(value)
	if dc.seen[cleaned] {
		return false
	}
	dc.seen[cleaned] = true
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
	
	samples := []string{"  Apple  ", "apple", "BANANA", "banana ", "  Cherry"}
	
	for _, sample := range samples {
		if cleaner.AddItem(sample) {
			fmt.Printf("Added: '%s'\n", sample)
		} else {
			fmt.Printf("Duplicate: '%s'\n", sample)
		}
	}
	
	fmt.Printf("Unique items: %d\n", cleaner.UniqueCount())
	
	cleaner.Reset()
	fmt.Println("Cleaner has been reset")
	fmt.Printf("Unique items after reset: %d\n", cleaner.UniqueCount())
}