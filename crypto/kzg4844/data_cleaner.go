
package main

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
	input := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(input)
	fmt.Printf("Original: %v\n", input)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[int]bool)
	var result []DataRecord
	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			result = append(result, record)
		}
	}
	return result
}

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	for _, record := range records {
		if ValidateEmail(record.Email) && record.Name != "" {
			cleaned = append(cleaned, record)
		}
	}
	return RemoveDuplicates(cleaned)
}

func main() {
	sampleData := []DataRecord{
		{1, "test@example.com", "John"},
		{2, "invalid-email", "Jane"},
		{3, "another@test.org", ""},
		{1, "test@example.com", "John"},
		{4, "valid@domain.com", "Alice"},
	}

	cleaned := CleanData(sampleData)
	fmt.Printf("Original records: %d\n", len(sampleData))
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Name: %s\n", record.ID, record.Email, record.Name)
	}
}