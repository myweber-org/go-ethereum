
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID      string
	Name    string
	Email   string
	Active  string
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

		if lineNumber == 1 {
			continue
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Name == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		if !strings.Contains(record.Email, "@") {
			return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) error {
	emailSet := make(map[string]bool)

	for _, record := range records {
		if emailSet[record.Email] {
			return fmt.Errorf("duplicate email found: %s", record.Email)
		}
		emailSet[record.Email] = true

		if record.Active != "true" && record.Active != "false" {
			return fmt.Errorf("invalid active status for record %s", record.ID)
		}
	}

	return nil
}

func GenerateReport(records []DataRecord) {
	activeCount := 0
	for _, record := range records {
		if record.Active == "true" {
			activeCount++
		}
	}

	fmt.Printf("Total records processed: %d\n", len(records))
	fmt.Printf("Active records: %d\n", activeCount)
	fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}