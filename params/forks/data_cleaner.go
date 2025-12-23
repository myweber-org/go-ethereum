
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Age   int
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("invalid record ID")
	}
	if record.Age < 0 || record.Age > 150 {
		return errors.New("invalid age value")
	}
	if !strings.Contains(record.Email, "@") {
		return errors.New("invalid email format")
	}
	
	dc.records = append(dc.records, record)
	return nil
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	unique := make([]DataRecord, 0)
	seen := make(map[int]bool)
	
	for _, record := range dc.records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}
	
	dc.records = unique
	return unique
}

func (dc *DataCleaner) ValidateAll() ([]DataRecord, []error) {
	validRecords := make([]DataRecord, 0)
	validationErrors := make([]error, 0)
	
	for _, record := range dc.records {
		err := dc.validateRecord(record)
		if err != nil {
			validationErrors = append(validationErrors, err)
		} else {
			validRecords = append(validRecords, record)
		}
	}
	
	return validRecords, validationErrors
}

func (dc *DataCleaner) validateRecord(record DataRecord) error {
	if record.Name == "" {
		return fmt.Errorf("record %d: name cannot be empty", record.ID)
	}
	if !strings.Contains(record.Email, "@") {
		return fmt.Errorf("record %d: invalid email format", record.ID)
	}
	if record.Age < 0 {
		return fmt.Errorf("record %d: age cannot be negative", record.ID)
	}
	return nil
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func (dc *DataCleaner) ClearAll() {
	dc.records = make([]DataRecord, 0)
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []DataRecord{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
		{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
		{ID: 3, Name: "", Email: "invalid-email", Age: -5},
	}
	
	for _, record := range records {
		err := cleaner.AddRecord(record)
		if err != nil {
			fmt.Printf("Error adding record %d: %v\n", record.ID, err)
		}
	}
	
	fmt.Printf("Total records before deduplication: %d\n", cleaner.GetRecordCount())
	
	unique := cleaner.RemoveDuplicates()
	fmt.Printf("Unique records: %d\n", len(unique))
	
	validRecords, errors := cleaner.ValidateAll()
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Validation errors: %d\n", len(errors))
	
	for _, err := range errors {
		fmt.Printf("Error: %v\n", err)
	}
}
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
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}