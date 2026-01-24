
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
	var validRecords []DataRecord

	for _, record := range records {
		cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
		if validateEmail(cleanEmail) && !emailMap[cleanEmail] {
			emailMap[cleanEmail] = true
			record.Valid = true
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func main() {
	emails := []string{
		"user@example.com",
		"USER@example.com",
		"user@example.com",
		"invalid-email",
		"another@test.org",
	}

	uniqueEmails := deduplicateEmails(emails)
	fmt.Println("Deduplicated emails:", uniqueEmails)

	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "admin@test.org", false},
		{3, "user@example.com", false},
		{4, "invalid", false},
	}

	processed := processRecords(records)
	fmt.Printf("Valid records: %d\n", len(processed))
	for _, r := range processed {
		fmt.Printf("ID: %d, Email: %s\n", r.ID, r.Email)
	}
}