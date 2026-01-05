
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

		record := Record{
			ID:    id,
			Name:  line[1],
			Value: value,
		}
		records = append(records, record)
		lineNum++
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}
		if record.Name == "" {
			return fmt.Errorf("record ID %d has empty name", record.ID)
		}
		if record.Value < 0 {
			return fmt.Errorf("record ID %d has negative value", record.ID)
		}
		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true
	}
	return nil
}

func CalculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}