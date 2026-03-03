
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
	result := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%d|%s|%s", record.ID, strings.ToLower(record.Email), strings.ToLower(record.Name))
		if !seen[key] {
			seen[key] = true
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateEmails() (valid []DataRecord, invalid []DataRecord) {
	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}
	return valid, invalid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(DataRecord{ID: 1, Email: "user@example.com", Name: "John Doe"})
	cleaner.AddRecord(DataRecord{ID: 2, Email: "user@example.com", Name: "John Doe"})
	cleaner.AddRecord(DataRecord{ID: 3, Email: "invalid-email", Name: "Jane Smith"})
	cleaner.AddRecord(DataRecord{ID: 4, Email: "another@test.org", Name: "Bob Wilson"})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", cleaner.GetRecordCount())

	valid, invalid := cleaner.ValidateEmails()
	fmt.Printf("Valid emails: %d, Invalid emails: %d\n", len(valid), len(invalid))
}
package utils

func DeduplicateStrings(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Value float64
}

type DataCleaner struct {
	records map[string]DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make(map[string]DataRecord),
	}
}

func (dc *DataCleaner) ValidateRecord(record DataRecord) error {
	if strings.TrimSpace(record.ID) == "" {
		return errors.New("empty ID")
	}
	if !strings.Contains(record.Email, "@") {
		return errors.New("invalid email format")
	}
	if record.Value < 0 {
		return errors.New("negative value not allowed")
	}
	return nil
}

func (dc *DataCleaner) AddRecord(record DataRecord) error {
	if err := dc.ValidateRecord(record); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if _, exists := dc.records[record.ID]; exists {
		return errors.New("duplicate ID")
	}
	dc.records[record.ID] = record
	return nil
}

func (dc *DataCleaner) RemoveDuplicatesByEmail() int {
	emailMap := make(map[string]string)
	removedCount := 0

	for id, record := range dc.records {
		if existingID, found := emailMap[record.Email]; found {
			delete(dc.records, id)
			removedCount++
			continue
		}
		emailMap[record.Email] = id
	}
	return removedCount
}

func (dc *DataCleaner) GetRecords() []DataRecord {
	records := make([]DataRecord, 0, len(dc.records))
	for _, record := range dc.records {
		records = append(records, record)
	}
	return records
}

func (dc *DataCleaner) CalculateAverage() float64 {
	if len(dc.records) == 0 {
		return 0
	}
	var sum float64
	for _, record := range dc.records {
		sum += record.Value
	}
	return sum / float64(len(dc.records))
}

func main() {
	cleaner := NewDataCleaner()

	sampleData := []DataRecord{
		{ID: "001", Email: "user1@example.com", Value: 10.5},
		{ID: "002", Email: "user2@example.com", Value: 20.3},
		{ID: "003", Email: "user1@example.com", Value: 15.7},
		{ID: "004", Email: "invalid-email", Value: 5.0},
	}

	for _, record := range sampleData {
		if err := cleaner.AddRecord(record); err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		}
	}

	fmt.Printf("Initial records: %d\n", len(cleaner.GetRecords()))
	removed := cleaner.RemoveDuplicatesByEmail()
	fmt.Printf("Removed duplicates: %d\n", removed)
	fmt.Printf("Final records: %d\n", len(cleaner.GetRecords()))
	fmt.Printf("Average value: %.2f\n", cleaner.CalculateAverage())
}