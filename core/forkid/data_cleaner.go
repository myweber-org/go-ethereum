
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
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    string
	Email string
	Score int
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
	if record.Score < 0 || record.Score > 100 {
		return errors.New("score must be between 0 and 100")
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
	validRecords := make([]DataRecord, 0, len(dc.records))
	for _, record := range dc.records {
		validRecords = append(validRecords, record)
	}
	return validRecords
}

func (dc *DataCleaner) CountRecords() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	sampleRecords := []DataRecord{
		{ID: "001", Email: "user1@example.com", Score: 85},
		{ID: "002", Email: "user2@example.com", Score: 92},
		{ID: "003", Email: "invalid-email", Score: 45},
		{ID: "001", Email: "duplicate@example.com", Score: 70},
		{ID: "004", Email: "user4@example.com", Score: 105},
	}

	for _, record := range sampleRecords {
		err := cleaner.AddRecord(record)
		if err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		} else {
			fmt.Printf("Added record %s successfully\n", record.ID)
		}
	}

	fmt.Printf("\nTotal valid records: %d\n", cleaner.CountRecords())
	
	validRecords := cleaner.GetValidRecords()
	fmt.Println("\nValid records:")
	for _, record := range validRecords {
		fmt.Printf("ID: %s, Email: %s, Score: %d\n", record.ID, record.Email, record.Score)
	}
}