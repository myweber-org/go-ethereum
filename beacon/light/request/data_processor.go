
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim spaces from edges
	cleaned = strings.TrimSpace(cleaned)
	
	// Convert to lowercase for normalization
	cleaned = strings.ToLower(cleaned)
	
	return cleaned
}

func NormalizeString(input string) string {
	cleaned := CleanInput(input)
	
	// Remove special characters except alphanumeric and spaces
	re := regexp.MustCompile(`[^a-z0-9\s]`)
	normalized := re.ReplaceAllString(cleaned, "")
	
	return normalized
}

func ProcessData(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		processed := NormalizeString(input)
		results = append(results, processed)
	}
	return results
}
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

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
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

	valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid boolean at line %d: %w", lineNumber, err)
	}
	record.Valid = valid

	return record, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []error) {
	var validRecords []DataRecord
	var validationErrors []error

	for i, record := range records {
		if record.ID <= 0 {
			validationErrors = append(validationErrors, 
				fmt.Errorf("record %d: ID must be positive, got %d", i+1, record.ID))
			continue
		}

		if len(record.Name) > 100 {
			validationErrors = append(validationErrors,
				fmt.Errorf("record %d: name exceeds 100 characters", i+1))
			continue
		}

		if record.Value < 0 {
			validationErrors = append(validationErrors,
				fmt.Errorf("record %d: value cannot be negative, got %f", i+1, record.Value))
			continue
		}

		validRecords = append(validRecords, record)
	}

	return validRecords, validationErrors
}

func CalculateStatistics(records []DataRecord) (float64, float64, error) {
	if len(records) == 0 {
		return 0, 0, errors.New("no records to calculate statistics")
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
		return 0, 0, errors.New("no valid records found")
	}

	average := sum / float64(count)

	var variance float64
	for _, record := range records {
		if record.Valid {
			diff := record.Value - average
			variance += diff * diff
		}
	}
	variance = variance / float64(count)

	return average, variance, nil
}

func ProcessDataFile(filename string) error {
	records, err := ParseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	validRecords, errors := ValidateRecords(records)
	if len(errors) > 0 {
		fmt.Printf("Validation errors (%d):\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	if len(validRecords) == 0 {
		return errors.New("no valid records after validation")
	}

	average, variance, err := CalculateStatistics(validRecords)
	if err != nil {
		return fmt.Errorf("statistics calculation failed: %w", err)
	}

	fmt.Printf("Processing complete:\n")
	fmt.Printf("  Total records: %d\n", len(records))
	fmt.Printf("  Valid records: %d\n", len(validRecords))
	fmt.Printf("  Average value: %.2f\n", average)
	fmt.Printf("  Variance: %.2f\n", variance)

	return nil
}