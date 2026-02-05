
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
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmailFormat(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func CleanData(records []DataRecord) []DataRecord {
	deduped := DeduplicateEmails(records)
	var cleaned []DataRecord

	for _, record := range deduped {
		record.Valid = ValidateEmailFormat(record.Email)
		cleaned = append(cleaned, record)
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

	cleaned := CleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
	}
}package main

import (
	"strings"
)

// DataCleaner provides methods for cleaning string data
type DataCleaner struct{}

// RemoveDuplicates removes duplicate entries from a slice of strings
func (dc *DataCleaner) RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// TrimSpaces removes leading and trailing whitespace from all strings in slice
func (dc *DataCleaner) TrimSpaces(input []string) []string {
	result := make([]string, len(input))
	for i, item := range input {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

// CleanData performs both duplicate removal and whitespace trimming
func (dc *DataCleaner) CleanData(input []string) []string {
	trimmed := dc.TrimSpaces(input)
	return dc.RemoveDuplicates(trimmed)
}