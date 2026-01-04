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
		if !seen[record.Email] {
			seen[record.Email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		record.Valid = strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".")
		validated = append(validated, record)
	}
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "test@domain.org", false},
	}

	uniqueData := RemoveDuplicates(sampleData)
	validatedData := ValidateEmails(uniqueData)

	for _, record := range validatedData {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}