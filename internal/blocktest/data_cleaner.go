package data

func DeduplicateStrings(slice []string) []string {
    seen := make(map[string]struct{})
    result := make([]string, 0, len(slice))
    
    for _, item := range slice {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    
    return result
}
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

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if strings.Contains(record.Email, "@") && len(record.Email) > 3 {
			record.Valid = true
		} else {
			record.Valid = false
		}
		valid = append(valid, record)
	}
	return valid
}

func PrintRecords(records []DataRecord) {
	for _, r := range records {
		status := "INVALID"
		if r.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", r.ID, r.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "  USER@EXAMPLE.COM  ", false},
		{2, "test@domain.org", false},
		{3, "user@example.com", false},
		{4, "invalid-email", false},
		{5, "another@test.com", false},
		{6, "TEST@DOMAIN.ORG", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	unique := DeduplicateRecords(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(unique)

	validated := ValidateEmails(unique)
	fmt.Println("\nAfter validation:")
	PrintRecords(validated)
}