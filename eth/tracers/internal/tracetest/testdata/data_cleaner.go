
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

func (dc *DataCleaner) AddRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if !strings.Contains(record.Email, "@") {
		return errors.New("invalid email format")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}

	if _, exists := dc.records[record.ID]; exists {
		return fmt.Errorf("duplicate record ID: %s", record.ID)
	}

	dc.records[record.ID] = record
	return nil
}

func (dc *DataCleaner) RemoveRecord(id string) bool {
	if _, exists := dc.records[id]; exists {
		delete(dc.records, id)
		return true
	}
	return false
}

func (dc *DataCleaner) GetValidRecords() []DataRecord {
	var validRecords []DataRecord
	for _, record := range dc.records {
		validRecords = append(validRecords, record)
	}
	return validRecords
}

func (dc *DataCleaner) Count() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	records := []DataRecord{
		{"user1", "test@example.com", 42.5},
		{"user2", "invalid-email", 100.0},
		{"user1", "another@example.com", 75.3},
		{"user3", "valid@domain.com", -10.0},
	}

	for _, record := range records {
		err := cleaner.AddRecord(record)
		if err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		}
	}

	fmt.Printf("Valid records count: %d\n", cleaner.Count())
	for _, record := range cleaner.GetValidRecords() {
		fmt.Printf("Record: %+v\n", record)
	}
}