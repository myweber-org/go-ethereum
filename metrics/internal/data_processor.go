
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
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("line %d: %v", lineNumber, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 columns, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row)
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNumber, err)
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID: %v", err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value: %v", err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return record, fmt.Errorf("invalid active flag: %v", err)
	}
	record.Active = active

	return record, nil
}

func validateRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.ID > 0 && record.Value >= 0 {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func calculateStats(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	for _, record := range records {
		sum += record.Value
		if record.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return sum, average, activeCount
}

func processDataFile(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return err
	}

	validRecords := validateRecords(records)
	total, average, activeCount := calculateStats(validRecords)

	fmt.Printf("Processed %d records (%d valid)\n", len(records), len(validRecords))
	fmt.Printf("Total value: %.2f\n", total)
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Active records: %d\n", activeCount)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	err := processDataFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}
}