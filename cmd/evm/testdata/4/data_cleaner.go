
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
	var valid []DataRecord
	for _, record := range records {
		if record.ID > 0 &&
			len(strings.TrimSpace(record.Name)) > 0 &&
			strings.Contains(record.Email, "@") {
			record.Valid = true
			valid = append(valid, record)
		}
	}
	return valid
}

func PrintRecords(records []DataRecord) {
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
		{4, "", "invalid-email", false},
		{5, "Bob Wilson", "bob@example.com", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	unique := RemoveDuplicates(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(unique)

	valid := ValidateRecords(unique)
	fmt.Println("\nAfter validation:")
	PrintRecords(valid)
}