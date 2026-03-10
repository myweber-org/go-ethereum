
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNumber)
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	validStr := strings.ToLower(strings.TrimSpace(row[3]))
	if validStr != "true" && validStr != "false" {
		return record, fmt.Errorf("invalid boolean at line %d: %s", lineNumber, row[3])
	}
	record.Valid = validStr == "true"

	return record, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []error) {
	var validRecords []DataRecord
	var validationErrors []error

	for _, record := range records {
		err := validateSingleRecord(record)
		if err != nil {
			validationErrors = append(validationErrors, err)
			continue
		}
		validRecords = append(validRecords, record)
	}

	return validRecords, validationErrors
}

func validateSingleRecord(record DataRecord) error {
	if record.ID <= 0 {
		return errors.New("ID must be positive")
	}

	if len(record.Name) > 100 {
		return errors.New("name exceeds maximum length")
	}

	if record.Value < 0 {
		return errors.New("value cannot be negative")
	}

	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64
	validCount := 0

	for i, record := range records {
		if !record.Valid {
			continue
		}

		sum += record.Value
		validCount++

		if i == 0 {
			min = record.Value
			max = record.Value
		} else {
			if record.Value < min {
				min = record.Value
			}
			if record.Value > max {
				max = record.Value
			}
		}
	}

	if validCount == 0 {
		return 0, 0, 0
	}

	average := sum / float64(validCount)
	return average, max - min, validCount
}