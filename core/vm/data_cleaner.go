
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedCount int
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{processedCount: 0}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; !exists {
			seen[trimmed] = struct{}{}
			result = append(result, trimmed)
			dc.processedCount++
		}
	}
	return result
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) GetStats() string {
	return fmt.Sprintf("Processed %d unique items", dc.processedCount)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"user@example.com",
		"  user@example.com  ",
		"invalid-email",
		"another@test.org",
		"",
		"another@test.org",
	}
	
	uniqueEmails := cleaner.RemoveDuplicates(data)
	fmt.Println("Unique emails:", uniqueEmails)
	
	for _, email := range uniqueEmails {
		if cleaner.ValidateEmail(email) {
			fmt.Printf("%s is valid\n", email)
		} else {
			fmt.Printf("%s is invalid\n", email)
		}
	}
	
	fmt.Println(cleaner.GetStats())
}