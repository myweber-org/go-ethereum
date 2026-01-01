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

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	unique := RemoveDuplicates(records)
	for _, record := range unique {
		if ValidateEmail(record.Email) {
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "test@example.com", false},
		{2, "invalid-email", false},
		{3, "test@example.com", false},
		{4, "user@domain.org", false},
	}
	cleaned := CleanData(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}