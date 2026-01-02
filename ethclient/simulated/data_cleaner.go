
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
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (dc *DataCleaner) TrimWhitespace() []string {
	result := make([]string, len(dc.Data))
	for i, item := range dc.Data {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func (dc *DataCleaner) Clean() []string {
	trimmed := dc.TrimWhitespace()
	dc.Data = trimmed
	return dc.RemoveDuplicates()
}

func main() {
	data := []string{"  apple ", "banana", "  apple", "cherry  ", "banana"}
	cleaner := NewDataCleaner(data)
	cleaned := cleaner.Clean()
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}