package csvutil

import (
	"strings"
)

func SanitizeCSVRow(row []string) []string {
	sanitized := make([]string, len(row))
	for i, field := range row {
		sanitized[i] = strings.TrimSpace(field)
	}
	return sanitized
}

func RemoveEmptyFields(row []string) []string {
	var result []string
	for _, field := range row {
		if strings.TrimSpace(field) != "" {
			result = append(result, field)
		}
	}
	return result
}