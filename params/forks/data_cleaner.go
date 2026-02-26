package csvutil

import (
	"strings"
)

func SanitizeCSVRow(row []string) []string {
	sanitized := make([]string, len(row))
	for i, cell := range row {
		sanitized[i] = strings.TrimSpace(cell)
	}
	return sanitized
}

func RemoveEmptyRows(rows [][]string) [][]string {
	var filtered [][]string
	for _, row := range rows {
		if !isEmptyRow(row) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}