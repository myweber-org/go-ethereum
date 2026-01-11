
package main

import (
	"strings"
)

type DataCleaner struct {
	trimSpaces bool
}

func NewDataCleaner(trimSpaces bool) *DataCleaner {
	return &DataCleaner{trimSpaces: trimSpaces}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range items {
		processed := item
		if dc.trimSpaces {
			processed = strings.TrimSpace(item)
		}
		
		if _, exists := seen[processed]; !exists && processed != "" {
			seen[processed] = struct{}{}
			result = append(result, processed)
		}
	}
	return result
}

func (dc *DataCleaner) CleanSlice(items []string) []string {
	cleaned := make([]string, 0, len(items))
	
	for _, item := range items {
		processed := item
		if dc.trimSpaces {
			processed = strings.TrimSpace(item)
		}
		if processed != "" {
			cleaned = append(cleaned, processed)
		}
	}
	return cleaned
}