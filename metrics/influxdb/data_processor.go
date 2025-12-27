
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID        int
	Name      string
	Value     float64
	IsActive  bool
	Timestamp string
}

func parseCSVFile(filename string) ([]DataRecord, error) {
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

		if len(row) != 5 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 5, got %d", lineNumber, len(row))
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
	var err error

	record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
	}

	record.Value, err = strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}

	record.IsActive, err = strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid boolean at line %d: %w", lineNumber, err)
	}

	record.Timestamp = strings.TrimSpace(row[4])
	if record.Timestamp == "" {
		return DataRecord{}, fmt.Errorf("empty timestamp at line %d", lineNumber)
	}

	return record, nil
}

func validateRecords(records []DataRecord) []error {
	var errors []error

	for i, record := range records {
		if record.ID <= 0 {
			errors = append(errors, fmt.Errorf("record %d: invalid ID %d", i+1, record.ID))
		}

		if record.Value < 0 {
			errors = append(errors, fmt.Errorf("record %d: negative value %f", i+1, record.Value))
		}

		if len(record.Name) > 100 {
			errors = append(errors, fmt.Errorf("record %d: name exceeds 100 characters", i+1))
		}
	}

	return errors
}

func processData(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	validationErrors := validateRecords(records)
	if len(validationErrors) > 0 {
		for _, err := range validationErrors {
			fmt.Printf("Validation error: %v\n", err)
		}
		return fmt.Errorf("validation found %d errors", len(validationErrors))
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	for _, record := range records {
		fmt.Printf("ID: %d, Name: %s, Value: %.2f, Active: %v, Time: %s\n",
			record.ID, record.Name, record.Value, record.IsActive, record.Timestamp)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	if err := processData(filename); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}