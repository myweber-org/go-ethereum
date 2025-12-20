
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d", lineNum)
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  line[1],
			Value: value,
		})
		lineNum++
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seen := make(map[int]bool)
	for _, rec := range records {
		if rec.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", rec.ID)
		}
		if rec.Name == "" {
			return fmt.Errorf("record ID %d has empty name", rec.ID)
		}
		if rec.Value < 0 {
			return fmt.Errorf("record ID %d has negative value", rec.ID)
		}
		if seen[rec.ID] {
			return fmt.Errorf("duplicate ID found: %d", rec.ID)
		}
		seen[rec.ID] = true
	}
	return nil
}

func CalculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, rec := range records {
		sum += rec.Value
		if rec.Value > max {
			max = rec.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}