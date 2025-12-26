
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
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformed := record
	transformed.Value = record.Value * 1.1

	if len(record.Tags) > 0 {
		transformed.Tags = make([]string, len(record.Tags))
		for i, tag := range record.Tags {
			transformed.Tags[i] = strings.ToUpper(strings.TrimSpace(tag))
		}
	}

	transformed.Timestamp = record.Timestamp.UTC()
	return transformed, nil
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	var errors []string

	for _, record := range records {
		transformed, err := TransformRecord(record)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Record %s: %v", record.ID, err))
			continue
		}
		processed = append(processed, transformed)
	}

	if len(errors) > 0 {
		return processed, fmt.Errorf("processing completed with errors: %v", strings.Join(errors, "; "))
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec1",
			Value:     100.0,
			Timestamp: time.Now(),
			Tags:      []string{"important", "test"},
		},
		{
			ID:        "rec2",
			Value:     -50.0,
			Timestamp: time.Now(),
			Tags:      []string{"invalid"},
		},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
	}

	fmt.Printf("Successfully processed %d records\n", len(processed))
	for _, rec := range processed {
		fmt.Printf("Record: %+v\n", rec)
	}
}