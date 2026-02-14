
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
	uniqueRecords := []string{}
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] && dc.isValidRecord(normalized) {
			dc.processedRecords[normalized] = true
			uniqueRecords = append(uniqueRecords, record)
		}
	}
	return uniqueRecords
}

func (dc *DataCleaner) isValidRecord(record string) bool {
	if record == "" {
		return false
	}
	if len(record) > 1000 {
		return false
	}
	return true
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
	return true
}

func (dc *DataCleaner) GetProcessedCount() int {
	return len(dc.processedRecords)
}

func main() {
	cleaner := NewDataCleaner()
	
	sampleData := []string{
		"user@example.com",
		"USER@EXAMPLE.COM",
		"test@domain.org",
		"invalid-email",
		"test@domain.org",
		"",
		"another@test.com",
	}
	
	fmt.Println("Original records:", len(sampleData))
	cleaned := cleaner.RemoveDuplicates(sampleData)
	fmt.Println("Cleaned records:", len(cleaned))
	fmt.Println("Processed count:", cleaner.GetProcessedCount())
	
	for _, email := range []string{"valid@test.com", "invalid"} {
		fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
	}
}