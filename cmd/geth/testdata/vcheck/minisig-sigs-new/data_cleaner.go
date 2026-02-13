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

func DeduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, DataRecord{
				ID:    record.ID,
				Email: email,
				Valid: record.Valid,
			})
		}
	}
	return unique
}

func ValidateEmailFormat(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}

func CleanDataset(records []DataRecord) []DataRecord {
	deduped := DeduplicateEmails(records)
	var cleaned []DataRecord

	for _, record := range deduped {
		isValid := ValidateEmailFormat(record.Email)
		cleaned = append(cleaned, DataRecord{
			ID:    record.ID,
			Email: record.Email,
			Valid: isValid,
		})
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain", false},
		{5, "user@example.com", false},
	}

	cleaned := CleanDataset(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned:  %d records\n", len(cleaned))

	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}