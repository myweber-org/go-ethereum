package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
	Category  string
}

func ValidateRecord(record DataRecord) error {
	var validationErrors []string

	if record.ID == "" {
		validationErrors = append(validationErrors, "ID cannot be empty")
	}

	if record.Timestamp.IsZero() {
		validationErrors = append(validationErrors, "Timestamp must be set")
	}

	if record.Value < 0 {
		validationErrors = append(validationErrors, "Value cannot be negative")
	}

	if record.Category == "" {
		validationErrors = append(validationErrors, "Category must be specified")
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}

	return nil
}

func TransformRecord(record DataRecord, multiplier float64) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, fmt.Errorf("validation failed: %w", err)
	}

	transformed := DataRecord{
		ID:        strings.ToUpper(record.ID),
		Timestamp: record.Timestamp.UTC(),
		Value:     record.Value * multiplier,
		Category:  strings.ToLower(record.Category),
	}

	return transformed, nil
}

func ProcessBatch(records []DataRecord, multiplier float64) ([]DataRecord, []error) {
	var processed []DataRecord
	var processingErrors []error

	for i, record := range records {
		transformed, err := TransformRecord(record, multiplier)
		if err != nil {
			processingErrors = append(processingErrors, fmt.Errorf("record %d: %w", i, err))
			continue
		}
		processed = append(processed, transformed)
	}

	return processed, processingErrors
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records provided")
	}

	var sum float64
	var count int

	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			continue
		}
		sum += record.Value
		count++
	}

	if count == 0 {
		return 0, 0, errors.New("no valid records found")
	}

	average := sum / float64(count)
	return sum, average, nil
}