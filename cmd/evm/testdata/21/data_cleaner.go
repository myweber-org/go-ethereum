package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	deduped := DeduplicateRecords(records)
	var cleaned []DataRecord
	for _, record := range deduped {
		record.Valid = ValidateEmail(record.Email)
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{ID: "001", Email: "user@example.com"},
		{ID: "002", Email: "invalid-email"},
		{ID: "001", Email: "user@example.com"},
		{ID: "003", Email: "another@test.org"},
	}

	cleaned := CleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %s, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}