
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func DeduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateRecord(record Record) error {
	if record.ID <= 0 {
		return errors.New("invalid record ID")
	}
	if strings.TrimSpace(record.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if !strings.Contains(record.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func CleanData(records []Record) ([]Record, error) {
	var cleaned []Record
	unique := DeduplicateRecords(records)

	for _, record := range unique {
		if err := ValidateRecord(record); err != nil {
			fmt.Printf("Skipping record %d: %v\n", record.ID, err)
			continue
		}
		record.Valid = true
		cleaned = append(cleaned, record)
	}

	if len(cleaned) == 0 {
		return nil, errors.New("no valid records after cleaning")
	}
	return cleaned, nil
}

func main() {
	sampleData := []Record{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "", "invalid-email", false},
		{5, "Alice Brown", "alice@example.com", false},
	}

	cleaned, err := CleanData(sampleData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Cleaned %d records from %d original\n", len(cleaned), len(sampleData))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}