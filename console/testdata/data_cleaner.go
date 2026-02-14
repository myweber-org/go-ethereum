package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
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

func DeduplicateEmails(emails []string) []string {
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

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return true
}

func CleanRecords(records []DataRecord) []DataRecord {
	emailSet := make(map[string]bool)
	var cleaned []DataRecord

	for _, record := range records {
		record.Email = strings.ToLower(strings.TrimSpace(record.Email))
		
		if !emailSet[record.Email] && ValidateEmail(record.Email) {
			emailSet[record.Email] = true
			record.Valid = true
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	emails := []string{
		"test@example.com",
		"TEST@example.com",
		"invalid-email",
		"another@test.org",
		"test@example.com",
	}

	uniqueEmails := DeduplicateEmails(emails)
	fmt.Printf("Unique emails: %v\n", uniqueEmails)

	records := []DataRecord{
		{1, "user1@domain.com", false},
		{2, "USER1@DOMAIN.COM", false},
		{3, "bad-email", false},
		{4, "user2@test.org", false},
	}

	cleaned := CleanRecords(records)
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	for _, r := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
	}
}