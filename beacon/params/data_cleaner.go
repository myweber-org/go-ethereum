package main

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
	data := []string{"apple", "banana", "apple", "cherry", "banana", "date"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	
	re := regexp.MustCompile(`\s+`)
	trimmed = re.ReplaceAllString(trimmed, " ")
	
	var result strings.Builder
	for _, r := range trimmed {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			result.WriteRune(r)
		}
	}
	
	return result.String()
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
		key := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[key] {
			seen[key] = true
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
		record.Valid = ValidateEmail(record.Email)
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "another@test.org", false},
		{5, "MIXED@EXAMPLE.COM", false},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Processed %d records, %d after cleaning\n", len(sampleData), len(cleaned))
	
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}