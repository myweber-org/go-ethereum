package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		normalizedEmail := normalizeEmail(record.Email)
		if !seen[normalizedEmail] {
			seen[normalizedEmail] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func printRecords(records []DataRecord) {
	for _, r := range records {
		fmt.Printf("ID: %d, Email: %s, Name: %s\n", r.ID, r.Email, r.Name)
	}
}

func main() {
	records := []DataRecord{
		{1, "  USER@Example.COM  ", "Alice"},
		{2, "user@example.com", "Bob"},
		{3, "test@domain.org", "Charlie"},
		{4, "TEST@DOMAIN.ORG ", "David"},
	}

	fmt.Println("Original records:")
	printRecords(records)

	uniqueRecords := deduplicateRecords(records)

	fmt.Println("\nDeduplicated records:")
	printRecords(uniqueRecords)
}