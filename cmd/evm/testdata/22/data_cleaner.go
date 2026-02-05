
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		key := fmt.Sprintf("%d|%s", record.ID, strings.ToLower(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if validateEmail(record.Email) && record.Name != "" {
			valid = append(valid, record)
		}
	}
	return deduplicateRecords(valid)
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", "John"},
		{2, "invalid-email", "Jane"},
		{3, "user@example.com", "John"},
		{4, "test@domain.com", ""},
		{5, "admin@site.org", "Admin"},
	}

	cleaned := cleanData(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Name: %s\n", r.ID, r.Email, r.Name)
	}
}