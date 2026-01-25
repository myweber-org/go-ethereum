
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseJSON(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if data.ID <= 0 {
		return nil, fmt.Errorf("invalid ID: must be positive integer")
	}
	if data.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if data.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	return &data, nil
}

func main() {
	jsonStr := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	parsedData, err := ValidateAndParseJSON([]byte(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed data: %+v\n", parsedData)
}
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Valid     bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	record.ID = strings.ToUpper(record.ID)
	record.Value = record.Value * 1.1
	return record
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{"abc123", 100.0, time.Now().Add(-time.Hour), true},
		{"def456", 200.0, time.Now().Add(-2 * time.Hour), true},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, record := range processed {
		fmt.Printf("Processed: ID=%s, Value=%.2f\n", record.ID, record.Value)
	}
}