
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(normalized) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && !strings.ContainsAny(item, "!@#$%")
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", "Apple", "banana", "", "cherry!", "banana", "date"}
	cleaned := cleaner.RemoveDuplicates(data)
	
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
	fmt.Printf("Unique count: %d\n", len(cleaned))
	
	cleaner.Reset()
	
	testData := []string{"test1", "test2", "test1"}
	fmt.Printf("Reset test: %v\n", cleaner.RemoveDuplicates(testData))
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
		record.Valid = strings.Contains(record.Email, "@") && len(record.Email) > 3
		validated = append(validated, record)
	}
	return validated
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
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "ANOTHER@EXAMPLE.COM", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	unique := RemoveDuplicates(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(unique)

	validated := ValidateEmails(unique)
	fmt.Println("\nAfter validation:")
	PrintRecords(validated)
}