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
	ID    int
	Name  string
	Value float64
	Valid bool
}

func processCSV(inputPath string) ([]Record, error) {
	file, err := os.Open(inputPath)
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

		if len(row) < 4 {
			continue
		}

		record := Record{}
		record.ID, _ = strconv.Atoi(strings.TrimSpace(row[0]))
		record.Name = strings.TrimSpace(row[1])
		record.Value, _ = strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		record.Valid = strings.ToLower(strings.TrimSpace(row[3])) == "true"

		if record.ID > 0 && record.Name != "" {
			records = append(records, record)
		}
	}

	return records, nil
}

func validateRecords(records []Record) []Record {
	var validRecords []Record
	for _, r := range records {
		if r.Valid && r.Value >= 0 {
			validRecords = append(validRecords, r)
		}
	}
	return validRecords
}

func calculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <input_file.csv>")
		return
	}

	records, err := processCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	validRecords := validateRecords(records)
	total := calculateTotal(validRecords)

	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Total value: %.2f\n", total)
}