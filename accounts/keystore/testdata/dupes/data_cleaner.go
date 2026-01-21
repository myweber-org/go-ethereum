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
		key := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	
	return strings.Contains(parts[1], ".")
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	
	for _, record := range records {
		if ValidateEmail(record.Email) {
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	
	return DeduplicateRecords(cleaned)
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "invalid-email", false},
		{3, "user@example.com", false},
		{4, "another@test.org", false},
		{5, "noatsign", false},
		{6, "valid@domain.co", false},
	}
	
	cleaned := CleanData(records)
	
	fmt.Printf("Original records: %d\n", len(records))
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", 
			record.ID, record.Email, record.Valid)
	}
}