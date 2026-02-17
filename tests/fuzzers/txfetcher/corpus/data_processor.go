
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
	ID        int
	Name      string
	Value     float64
	Valid     bool
	Timestamp string
}

func ProcessCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
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

		if len(row) < 5 {
			continue
		}

		record, err := parseRecord(row)
		if err != nil {
			fmt.Printf("Skipping invalid record at line %d: %v\n", lineNumber, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRecord(row []string) (Record, error) {
	var record Record
	var err error

	record.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID: %w", err)
	}

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return Record{}, fmt.Errorf("name cannot be empty")
	}

	record.Value, err = strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value: %w", err)
	}

	record.Valid, err = strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid boolean flag: %w", err)
	}

	record.Timestamp = strings.TrimSpace(row[4])
	if record.Timestamp == "" {
		return Record{}, fmt.Errorf("timestamp cannot be empty")
	}

	return record, nil
}

func FilterValidRecords(records []Record) []Record {
	var validRecords []Record
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverageValue(records []Record) float64 {
	if len(records) == 0 {
		return 0.0
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}

func GenerateReport(records []Record) {
	validRecords := FilterValidRecords(records)
	avgValue := CalculateAverageValue(validRecords)

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Average value of valid records: %.2f\n", avgValue)

	if len(validRecords) > 0 {
		fmt.Println("\nValid Records Summary:")
		for _, record := range validRecords {
			fmt.Printf("ID: %d, Name: %s, Value: %.2f\n", record.ID, record.Name, record.Value)
		}
	}
}