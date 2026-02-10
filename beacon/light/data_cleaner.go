package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
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
	Email   string
	Active  bool
	Score   float64
}

func cleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	headers = append(headers, "Validated")
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	lineNum := 1
	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Warning: line %d: %v\n", lineNum, err)
			continue
		}

		record, validationErr := validateRow(row)
		validated := "PASS"
		if validationErr != nil {
			validated = "FAIL: " + validationErr.Error()
		}

		outputRow := []string{
			strconv.Itoa(record.ID),
			strings.TrimSpace(record.Name),
			strings.ToLower(strings.TrimSpace(record.Email)),
			strconv.FormatBool(record.Active),
			fmt.Sprintf("%.2f", record.Score),
			validated,
		}

		if err := writer.Write(outputRow); err != nil {
			fmt.Printf("Warning: failed to write line %d: %v\n", lineNum, err)
		}
	}

	return nil
}

func validateRow(row []string) (Record, error) {
	if len(row) < 5 {
		return Record{}, fmt.Errorf("insufficient columns")
	}

	var r Record
	var err error

	if r.ID, err = strconv.Atoi(row[0]); err != nil {
		return Record{}, fmt.Errorf("invalid ID: %w", err)
	}

	r.Name = row[1]
	if r.Name == "" {
		return Record{}, fmt.Errorf("name cannot be empty")
	}

	r.Email = row[2]
	if !strings.Contains(r.Email, "@") {
		return Record{}, fmt.Errorf("invalid email format")
	}

	if r.Active, err = strconv.ParseBool(row[3]); err != nil {
		return Record{}, fmt.Errorf("invalid active flag: %w", err)
	}

	if r.Score, err = strconv.ParseFloat(row[4], 64); err != nil {
		return Record{}, fmt.Errorf("invalid score: %w", err)
	}

	if r.Score < 0 || r.Score > 100 {
		return Record{}, fmt.Errorf("score out of range (0-100)")
	}

	return r, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data cleaning completed. Output saved to %s\n", outputFile)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	records []string
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]string, 0),
	}
}

func (dc *DataCleaner) AddRecord(record string) {
	dc.records = append(dc.records, strings.TrimSpace(record))
}

func (dc *DataCleaner) RemoveDuplicates() []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, record := range dc.records {
		if !seen[record] {
			seen[record] = true
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateRecords() (valid []string, invalid []string) {
	valid = make([]string, 0)
	invalid = make([]string, 0)

	for _, record := range dc.records {
		if len(record) > 0 && !strings.ContainsAny(record, "!@#$%") {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}

	return valid, invalid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	sampleData := []string{
		"user@example.com",
		"  user@example.com  ",
		"test@domain.org",
		"",
		"invalid!@email.com",
		"another@test.com",
		"test@domain.org",
	}

	for _, data := range sampleData {
		cleaner.AddRecord(data)
	}

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	deduped := cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", len(deduped))

	valid, invalid := cleaner.ValidateRecords()
	fmt.Printf("Valid records: %d\n", len(valid))
	fmt.Printf("Invalid records: %d\n", len(invalid))

	fmt.Println("\nValid records:")
	for _, v := range valid {
		fmt.Printf("  - %s\n", v)
	}
}