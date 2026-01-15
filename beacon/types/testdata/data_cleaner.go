
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	for i := range records {
		email := records[i].Email
		records[i].Valid = strings.Contains(email, "@") &&
			strings.Contains(email, ".") &&
			len(email) > 5
	}
	return records
}

func CleanData(records []DataRecord) []DataRecord {
	records = RemoveDuplicates(records)
	records = ValidateEmails(records)
	return records
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "ANOTHER@TEST.ORG", false},
	}

	cleaned := CleanData(sampleData)

	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n",
			record.ID, record.Email, record.Valid)
	}
}