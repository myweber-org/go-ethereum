
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
package main

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

func (dc *DataCleaner) Normalize(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func (dc *DataCleaner) IsDuplicate(value string) bool {
	normalized := dc.Normalize(value)
	if dc.seen[normalized] {
		return true
	}
	dc.seen[normalized] = true
	return false
}

func (dc *DataCleaner) Deduplicate(values []string) []string {
	dc.seen = make(map[string]bool)
	var result []string
	for _, v := range values {
		if !dc.IsDuplicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{"apple", " Apple ", "BANANA", "banana", "Cherry", "cherry "}
	
	fmt.Println("Original data:", data)
	
	deduped := cleaner.Deduplicate(data)
	fmt.Println("Deduplicated data:", deduped)
	
	testValue := "  APPLE  "
	if cleaner.IsDuplicate(testValue) {
		fmt.Printf("'%s' is a duplicate\n", testValue)
	} else {
		fmt.Printf("'%s' is unique\n", testValue)
	}
}