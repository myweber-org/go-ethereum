
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

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) Deduplicate(values []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, v := range values {
		if !dc.IsDuplicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"Apple", "apple", " BANANA ", "banana", "Cherry"}
	fmt.Println("Original:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated:", deduped)
	
	cleaner.Reset()
	testValue := "  TEST  "
	fmt.Printf("Normalized '%s': '%s'\n", testValue, cleaner.Normalize(testValue))
}package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
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
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanRecords(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    for _, record := range records {
        email := strings.ToLower(strings.TrimSpace(record.Email))
        if !emailSet[email] && ValidateEmail(email) {
            emailSet[email] = true
            record.Email = email
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
        {3, "invalid-email", false},
        {4, "test@domain.org", false},
    }
    
    cleaned := CleanRecords(records)
    fmt.Printf("Cleaned %d records\n", len(cleaned))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
}