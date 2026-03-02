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
	ID        int
	Name      string
	Email     string
	Age       int
	Validated bool
}

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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

		if len(row) < 4 {
			continue
		}

		id, _ := strconv.Atoi(strings.TrimSpace(row[0]))
		name := strings.TrimSpace(row[1])
		email := strings.TrimSpace(row[2])
		age, _ := strconv.Atoi(strings.TrimSpace(row[3]))

		record := DataRecord{
			ID:    id,
			Name:  name,
			Email: email,
			Age:   age,
		}

		record.Validated = validateRecord(record)
		records = append(records, record)
	}

	return records, nil
}

func validateRecord(record DataRecord) bool {
	if record.ID <= 0 {
		return false
	}
	if len(record.Name) == 0 {
		return false
	}
	if !strings.Contains(record.Email, "@") {
		return false
	}
	if record.Age < 0 || record.Age > 120 {
		return false
	}
	return true
}

func filterValidRecords(records []DataRecord) []DataRecord {
	valid := []DataRecord{}
	for _, record := range records {
		if record.Validated {
			valid = append(valid, record)
		}
	}
	return valid
}

func generateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	validCount := 0
	for _, record := range records {
		if record.Validated {
			validCount++
		}
	}
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Invalid records: %d\n", len(records)-validCount)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <csv_file>")
		return
	}

	filename := os.Args[1]
	records, err := parseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}

	generateReport(records)
	validRecords := filterValidRecords(records)
	fmt.Printf("Filtered valid records: %d\n", len(validRecords))
}