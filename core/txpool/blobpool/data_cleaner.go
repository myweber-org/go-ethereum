
package main

import (
	"strings"
)

type DataCleaner struct {
	trimSpaces bool
}

func NewDataCleaner(trimSpaces bool) *DataCleaner {
	return &DataCleaner{trimSpaces: trimSpaces}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range items {
		processed := item
		if dc.trimSpaces {
			processed = strings.TrimSpace(item)
		}
		
		if _, exists := seen[processed]; !exists && processed != "" {
			seen[processed] = struct{}{}
			result = append(result, processed)
		}
	}
	return result
}

func (dc *DataCleaner) CleanSlice(items []string) []string {
	cleaned := make([]string, 0, len(items))
	
	for _, item := range items {
		processed := item
		if dc.trimSpaces {
			processed = strings.TrimSpace(item)
		}
		if processed != "" {
			cleaned = append(cleaned, processed)
		}
	}
	return cleaned
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

func deduplicateEmails(emails []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if !seen[email] {
			seen[email] = true
			result = append(result, email)
		}
	}
	return result
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processRecords(records []DataRecord) []DataRecord {
	emailMap := make(map[string]bool)
	var cleaned []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(cleanEmail) && !emailMap[cleanEmail] {
			emailMap[cleanEmail] = true
			record.Email = cleanEmail
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
		{5, "test@domain.org", false},
	}

	cleaned := processRecords(records)
	fmt.Printf("Processed %d records, %d valid unique records found\n", len(records), len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}