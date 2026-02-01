package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNum, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 columns, got %d", lineNum, len(row))
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid ID: %v", lineNum, err)
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			return nil, fmt.Errorf("line %d: name cannot be empty", lineNum)
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid value: %v", lineNum, err)
		}

		active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid active flag: %v", lineNum, err)
		}

		records = append(records, Record{
			ID:     id,
			Name:   name,
			Value:  value,
			Active: active,
		})
	}

	return records, nil
}

func calculateStats(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var activeCount int
	var minVal float64 = records[0].Value

	for _, r := range records {
		sum += r.Value
		if r.Value < minVal {
			minVal = r.Value
		}
		if r.Active {
			activeCount++
		}
	}

	average := sum / float64(len(records))
	return average, minVal, activeCount
}

func filterRecords(records []Record, predicate func(Record) bool) []Record {
	var filtered []Record
	for _, r := range records {
		if predicate(r) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := parseCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully loaded %d records\n", len(records))

	avg, min, active := calculateStats(records)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Minimum value: %.2f\n", min)
	fmt.Printf("Active records: %d\n", active)

	activeRecords := filterRecords(records, func(r Record) bool {
		return r.Active && r.Value > 50.0
	})
	fmt.Printf("Records with value > 50 and active: %d\n", len(activeRecords))

	for i, r := range activeRecords {
		if i < 3 {
			fmt.Printf("  %d: %s (%.2f)\n", r.ID, r.Name, r.Value)
		}
	}
}