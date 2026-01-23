
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
	for i := range records {
		records[i].Valid = strings.Contains(records[i].Email, "@") && strings.Contains(records[i].Email, ".")
	}
	return records
}

func PrintRecords(records []DataRecord) {
	for _, record := range records {
		status := "Invalid"
		if record.Valid {
			status = "Valid"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", record.ID, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "another@domain.org", false},
		{5, "no.at.sign", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	uniqueRecords := RemoveDuplicates(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(uniqueRecords)

	validatedRecords := ValidateEmails(uniqueRecords)
	fmt.Println("\nAfter validation:")
	PrintRecords(validatedRecords)
}