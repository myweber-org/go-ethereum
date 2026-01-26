
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
	if !isValidEmail(record.Email) {
		return errors.New("invalid email format")
	}
	if record.Score < 0 || record.Score > 100 {
		return errors.New("score must be between 0 and 100")
	}

	if _, exists := dc.records[record.ID]; exists {
		return errors.New("duplicate record ID")
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

func (dc *DataCleaner) GetRecord(id string) (DataRecord, bool) {
	record, exists := dc.records[id]
	return record, exists
}

func (dc *DataCleaner) GetAllRecords() []DataRecord {
	records := make([]DataRecord, 0, len(dc.records))
	for _, record := range dc.records {
		records = append(records, record)
	}
	return records
}

func (dc *DataCleaner) CountRecords() int {
	return len(dc.records)
}

func (dc *DataCleaner) CalculateAverageScore() float64 {
	if len(dc.records) == 0 {
		return 0.0
	}

	total := 0
	for _, record := range dc.records {
		total += record.Score
	}
	return float64(total) / float64(len(dc.records))
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func main() {
	cleaner := NewDataCleaner()

	records := []DataRecord{
		{ID: "001", Email: "user1@example.com", Score: 85},
		{ID: "002", Email: "user2@example.org", Score: 92},
		{ID: "003", Email: "user3@example.net", Score: 78},
	}

	for _, record := range records {
		if err := cleaner.AddRecord(record); err != nil {
			fmt.Printf("Failed to add record %s: %v\n", record.ID, err)
		}
	}

	fmt.Printf("Total records: %d\n", cleaner.CountRecords())
	fmt.Printf("Average score: %.2f\n", cleaner.CalculateAverageScore())

	allRecords := cleaner.GetAllRecords()
	for _, record := range allRecords {
		fmt.Printf("ID: %s, Email: %s, Score: %d\n", record.ID, record.Email, record.Score)
	}
}