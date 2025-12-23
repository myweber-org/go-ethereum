
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range dc.Data {
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

func (dc *DataCleaner) Clean() []string {
	cleaned := dc.RemoveDuplicates()
	var final []string
	for _, item := range cleaned {
		final = append(final, strings.ToLower(item))
	}
	return final
}

func main() {
	rawData := []string{"  Apple ", "banana", "  Apple", "Banana", "  ", "cherry"}
	cleaner := NewDataCleaner(rawData)
	cleanedData := cleaner.Clean()
	fmt.Println("Cleaned data:", cleanedData)
}