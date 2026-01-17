
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
		if !seen[email] {
			seen[email] = true
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

func PrintRecords(records []DataRecord) {
	for _, record := range records {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", record.ID, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "test@domain.org", false},
		{3, "user@example.com", false},
		{4, "invalid-email", false},
		{5, "another@test.com", false},
		{6, "test@domain.org", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	uniqueRecords := RemoveDuplicates(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(uniqueRecords)

	validatedRecords := ValidateEmails(uniqueRecords)
	fmt.Println("\nAfter validation:")
	PrintRecords(validatedRecords)
}package utils

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}