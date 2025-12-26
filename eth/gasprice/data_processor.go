package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVData(reader io.Reader) ([]DataRecord, error) {
	csvReader := csv.NewReader(reader)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	var data []DataRecord
	for i, row := range records {
		if len(row) < 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			continue
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			continue
		}

		valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

		data = append(data, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		})
	}

	return data, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, []DataRecord) {
	var valid []DataRecord
	var invalid []DataRecord

	for _, record := range records {
		if record.ID > 0 && record.Value >= 0 {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}

	return valid, invalid
}

func CalculateStats(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var count int
	min := records[0].Value

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			count++
			if record.Value < min {
				min = record.Value
			}
		}
	}

	if count == 0 {
		return 0, 0, 0
	}

	average := sum / float64(count)
	return average, min, count
}