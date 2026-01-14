package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateRecords(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		record.Valid = record.ID > 0 &&
			len(strings.TrimSpace(record.Name)) > 0 &&
			strings.Contains(record.Email, "@")
		validated = append(validated, record)
	}
	return validated
}

func PrintRecords(records []DataRecord) {
	fmt.Println("Processed Records:")
	for _, record := range records {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Status: %s\n",
			record.ID, record.Name, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{0, "", "invalid-email", false},
		{4, "Alice Brown", "alice@example.com", false},
	}

	fmt.Println("Original records count:", len(records))
	uniqueRecords := RemoveDuplicates(records)
	fmt.Println("After deduplication:", len(uniqueRecords))

	validatedRecords := ValidateRecords(uniqueRecords)
	PrintRecords(validatedRecords)
}