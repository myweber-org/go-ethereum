
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
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
	records := []DataRecord{}
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

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]
		if name == "" {
			continue
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		valid := false
		if row[3] == "true" {
			valid = true
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		}
		records = append(records, record)
	}

	return records, nil
}

func CalculateStats(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var validCount int
	var maxValue float64

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			validCount++
			if record.Value > maxValue {
				maxValue = record.Value
			}
		}
	}

	average := 0.0
	if validCount > 0 {
		average = sum / float64(validCount)
	}

	return average, maxValue, validCount
}