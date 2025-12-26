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