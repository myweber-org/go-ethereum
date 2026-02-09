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
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID   string
	Name string
	Email string
}

func generateHash(r Record) string {
	data := fmt.Sprintf("%s-%s", r.Name, r.Email)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func deduplicateRecords(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, record := range records {
		hash := generateHash(record)
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanData(records []Record) []Record {
	var valid []Record
	for _, record := range records {
		if validateEmail(record.Email) {
			valid = append(valid, record)
		}
	}
	return deduplicateRecords(valid)
}

func main() {
	sampleData := []Record{
		{ID: "1", Name: "John Doe", Email: "john@example.com"},
		{ID: "2", Name: "Jane Smith", Email: "jane@example.org"},
		{ID: "3", Name: "John Doe", Email: "john@example.com"},
		{ID: "4", Name: "Bob", Email: "invalid-email"},
	}

	cleaned := cleanData(sampleData)
	fmt.Printf("Original: %d records\n", len(sampleData))
	fmt.Printf("Cleaned: %d records\n", len(cleaned))
}package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processed map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processed: make(map[string]bool),
	}
}

func (dc *DataCleaner) Clean(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("empty input")
	}

	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)

	if dc.processed[lower] {
		return "", fmt.Errorf("duplicate entry: %s", trimmed)
	}

	dc.processed[lower] = true
	return trimmed, nil
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func main() {
	cleaner := NewDataCleaner()

	samples := []string{
		"  User@Example.com  ",
		"user@example.com",
		"test",
		"",
	}

	for _, sample := range samples {
		cleaned, err := cleaner.Clean(sample)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		isValid := cleaner.ValidateEmail(cleaned)
		fmt.Printf("Original: %q -> Cleaned: %q (Valid email: %v)\n", sample, cleaned, isValid)
	}
}