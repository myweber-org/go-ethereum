
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	Data []string
}

func NewDataCleaner(data []string) *DataCleaner {
	return &DataCleaner{Data: data}
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range dc.Data {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func (dc *DataCleaner) TrimWhitespace() []string {
	result := make([]string, len(dc.Data))
	for i, item := range dc.Data {
		result[i] = strings.TrimSpace(item)
	}
	return result
}

func (dc *DataCleaner) Clean() []string {
	trimmed := dc.TrimWhitespace()
	dc.Data = trimmed
	return dc.RemoveDuplicates()
}

func main() {
	data := []string{"  apple ", "banana", "  apple", "cherry  ", "banana"}
	cleaner := NewDataCleaner(data)
	cleaned := cleaner.Clean()
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var unique []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) SanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	replacer := strings.NewReplacer("\n", " ", "\r", " ", "\t", " ")
	return replacer.Replace(trimmed)
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{"user@example.com", "  USER@EXAMPLE.COM  ", "test@domain.org", "invalid-email"}
	
	fmt.Println("Original records:", records)
	
	uniqueRecords := cleaner.RemoveDuplicates(records)
	fmt.Println("After deduplication:", uniqueRecords)
	
	for _, record := range uniqueRecords {
		sanitized := cleaner.SanitizeInput(record)
		isValid := cleaner.ValidateEmail(sanitized)
		fmt.Printf("Email: %s, Valid: %v\n", sanitized, isValid)
	}
}