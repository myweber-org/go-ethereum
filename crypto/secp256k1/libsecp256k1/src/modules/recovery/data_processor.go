
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

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
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
			return nil, fmt.Errorf("invalid column count at line %d", lineNumber)
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

		records = append(records, record)
	}

	return records, nil
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func FilterActiveUsers(records []DataRecord) []DataRecord {
	var activeUsers []DataRecord
	for _, record := range records {
		if strings.ToLower(record.Active) == "true" {
			activeUsers = append(activeUsers, record)
		}
	}
	return activeUsers
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	
	activeUsers := FilterActiveUsers(records)
	fmt.Printf("Active users: %d\n", len(activeUsers))

	emailValidCount := 0
	for _, record := range records {
		if ValidateEmail(record.Email) {
			emailValidCount++
		}
	}
	fmt.Printf("Valid email addresses: %d\n", emailValidCount)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file_path>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	GenerateReport(records)
}