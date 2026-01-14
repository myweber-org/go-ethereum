
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", strings.ToLower(record.Name), strings.ToLower(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return len(email) > 5 && len(email) < 255
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	unique := DeduplicateRecords(records)

	for _, record := range unique {
		record.Valid = ValidateEmail(record.Email)
		if record.Valid {
			record.Name = strings.TrimSpace(record.Name)
			record.Email = strings.ToLower(strings.TrimSpace(record.Email))
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@test.org", false},
		{3, "John Doe", "JOHN@example.com", false},
		{4, "Bob", "invalid-email", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d valid records\n", len(cleaned))
	
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", r.ID, r.Name, r.Email)
	}
}