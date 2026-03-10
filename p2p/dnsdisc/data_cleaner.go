package datautils

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	slice := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	uniqueSlice := RemoveDuplicates(slice)
	fmt.Println("Original:", slice)
	fmt.Println("Unique:", uniqueSlice)
}package main

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
	
	data := []string{
		"apple",
		"Apple",
		"banana",
		"banana ",
		"",
		"cherry!",
		"date",
	}
	
	cleaned := cleaner.RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
	
	cleaner.Reset()
	fmt.Println("Cleaner has been reset")
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateEmails(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) error {
	if len(email) == 0 {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @ symbol")
	}
	if !strings.Contains(email, ".") {
		return errors.New("email must contain domain")
	}
	return nil
}

func CleanRecords(records []Record) ([]Record, error) {
	var cleaned []Record
	for _, record := range records {
		if err := ValidateEmail(record.Email); err != nil {
			record.Valid = false
		} else {
			record.Valid = true
		}
		cleaned = append(cleaned, record)
	}
	return DeduplicateEmails(cleaned), nil
}

func main() {
	records := []Record{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "test@domain.org", false},
		{4, "invalid-email", false},
	}

	cleaned, err := CleanRecords(records)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, r := range cleaned {
		status := "valid"
		if !r.Valid {
			status = "invalid"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", r.ID, r.Email, status)
	}
}