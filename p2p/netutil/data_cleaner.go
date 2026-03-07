
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

func (dc *DataCleaner) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := strings.ToLower(trimmed)
	return normalized
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	cleaned := dc.CleanString(value)
	if dc.seen[cleaned] {
		return true
	}
	dc.seen[cleaned] = true
	return false
}

func (dc *DataCleaner) AddItem(value string) bool {
	cleaned := dc.CleanString(value)
	if dc.seen[cleaned] {
		return false
	}
	dc.seen[cleaned] = true
	return true
}

func (dc *DataCleaner) UniqueCount() int {
	return len(dc.seen)
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	samples := []string{"  Apple  ", "apple", "BANANA", "banana ", "  Cherry"}
	
	for _, sample := range samples {
		if cleaner.AddItem(sample) {
			fmt.Printf("Added: '%s'\n", sample)
		} else {
			fmt.Printf("Duplicate: '%s'\n", sample)
		}
	}
	
	fmt.Printf("Unique items: %d\n", cleaner.UniqueCount())
	
	cleaner.Reset()
	fmt.Println("Cleaner has been reset")
	fmt.Printf("Unique items after reset: %d\n", cleaner.UniqueCount())
}
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

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
		if !seen[email] && email != "" {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	unique := DeduplicateRecords(records)

	for _, record := range unique {
		record.Email = strings.ToLower(strings.TrimSpace(record.Email))
		record.Valid = ValidateEmail(record.Email)
		if record.Valid {
			cleaned = append(cleaned, record)
		}
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain", false},
		{5, "another@test.org", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d valid records\n", len(cleaned))

	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s\n", record.ID, record.Email)
	}
}