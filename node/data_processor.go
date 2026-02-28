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
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
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
			continue
		}

		record, err := parseRow(row)
		if err != nil {
			fmt.Printf("skipping invalid row at line %d: %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID: %w", err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, fmt.Errorf("empty name")
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value: %w", err)
	}
	record.Value = value

	valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"
	record.Valid = valid

	return record, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var validCount int

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			validCount++
		}
	}

	if validCount == 0 {
		return 0, 0, 0
	}

	average := sum / float64(validCount)
	return sum, average, validCount
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Valid {
			filtered = append(filtered, record)
		}
	}
	return filtered
}