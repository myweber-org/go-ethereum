
package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateAndNormalize(data UserData) (UserData, error) {
	var normalized UserData

	normalized.Username = strings.TrimSpace(data.Username)
	if normalized.Username == "" {
		return UserData{}, fmt.Errorf("username cannot be empty")
	}

	normalized.Email = strings.ToLower(strings.TrimSpace(data.Email))
	if !strings.Contains(normalized.Email, "@") {
		return UserData{}, fmt.Errorf("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return UserData{}, fmt.Errorf("age must be between 0 and 150")
	}
	normalized.Age = data.Age

	return normalized, nil
}

func main() {
	testData := UserData{
		Username: "  JohnDoe  ",
		Email:    "  TEST@EXAMPLE.COM",
		Age:      25,
	}

	result, err := ValidateAndNormalize(testData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Normalized Data: %+v\n", result)
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
	}
}

func (dp *DataProcessor) Process() error {
	inputFile, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dp.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	cleanedCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		recordCount++
		cleanedRecord := dp.cleanRecord(record)

		if dp.isValidRecord(cleanedRecord) {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			cleanedCount++
		}
	}

	fmt.Printf("Processing complete. Total records: %d, Cleaned records: %d\n", recordCount, cleanedCount)
	return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
	}
	return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if field == "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}