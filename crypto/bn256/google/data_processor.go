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
	Value   string
	IsValid bool
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
	headerSkipped := false

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Value: strings.TrimSpace(row[2]),
		}
		record.IsValid = validateRecord(record)

		records = append(records, record)
	}

	return records, nil
}

func validateRecord(record DataRecord) bool {
	if record.ID == "" || record.Name == "" {
		return false
	}
	if len(record.Value) > 100 {
		return false
	}
	return true
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if record.IsValid {
			valid = append(valid, record)
		}
	}
	return valid
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	validRecords := FilterValidRecords(records)
	fmt.Printf("Total records: %d, Valid records: %d\n", len(records), len(validRecords))

	for _, record := range validRecords {
		fmt.Printf("ID: %s, Name: %s, Value: %s\n", record.ID, record.Name, record.Value)
	}
}