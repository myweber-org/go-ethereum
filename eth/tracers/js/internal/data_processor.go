
package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Timestamp time.Time
	Value     float64
	Category  string
	Valid     bool
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}
	if record.Category == "" {
		return errors.New("category cannot be empty")
	}
	return nil
}

func TransformCategory(input string) string {
	return strings.ToUpper(strings.TrimSpace(input))
}

func CalculateAverage(records []DataRecord) (float64, error) {
	if len(records) == 0 {
		return 0, errors.New("no records provided")
	}

	var sum float64
	var count int

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			count++
		}
	}

	if count == 0 {
		return 0, errors.New("no valid records found")
	}

	return sum / float64(count), nil
}

func FilterByCategory(records []DataRecord, category string) []DataRecord {
	var filtered []DataRecord
	targetCategory := TransformCategory(category)

	for _, record := range records {
		if TransformCategory(record.Category) == targetCategory && record.Valid {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func ProcessRecords(records []DataRecord) (map[string]float64, error) {
	results := make(map[string]float64)

	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			continue
		}

		record.Category = TransformCategory(record.Category)
		record.Valid = true

		categoryRecords := FilterByCategory(records, record.Category)
		average, err := CalculateAverage(categoryRecords)
		if err != nil {
			continue
		}

		results[record.Category] = average
	}

	if len(results) == 0 {
		return nil, errors.New("no valid data processed")
	}

	return results, nil
}