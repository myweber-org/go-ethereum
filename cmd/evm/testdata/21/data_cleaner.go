package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	deduped := DeduplicateRecords(records)
	var cleaned []DataRecord
	for _, record := range deduped {
		record.Valid = ValidateEmail(record.Email)
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{ID: "001", Email: "user@example.com"},
		{ID: "002", Email: "invalid-email"},
		{ID: "001", Email: "user@example.com"},
		{ID: "003", Email: "another@test.org"},
	}

	cleaned := CleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %s, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
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

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]bool)
	unique := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s", record.ID, strings.ToLower(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}

	dc.records = unique
	return unique
}

func (dc *DataCleaner) ValidateEmails() []DataRecord {
	valid := make([]DataRecord, 0)

	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && len(record.Name) > 0 {
			valid = append(valid, record)
		}
	}

	return valid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(DataRecord{ID: 1, Email: "user@example.com", Name: "Alice"})
	cleaner.AddRecord(DataRecord{ID: 2, Email: "user@example.com", Name: "Bob"})
	cleaner.AddRecord(DataRecord{ID: 3, Email: "invalid-email", Name: "Charlie"})
	cleaner.AddRecord(DataRecord{ID: 1, Email: "user@example.com", Name: "Alice"})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	unique := cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", len(unique))

	valid := cleaner.ValidateEmails()
	fmt.Printf("Valid records: %d\n", len(valid))

	for _, record := range valid {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}