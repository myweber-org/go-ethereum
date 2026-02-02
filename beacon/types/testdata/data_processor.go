
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) (DataRecord, error) {
	if err := ValidateRecord(record); err != nil {
		return DataRecord{}, err
	}

	transformed := record
	transformed.Value = record.Value * 1.1

	if len(record.Tags) > 0 {
		transformed.Tags = make([]string, len(record.Tags))
		for i, tag := range record.Tags {
			transformed.Tags[i] = strings.ToUpper(strings.TrimSpace(tag))
		}
	}

	transformed.Timestamp = record.Timestamp.UTC()
	return transformed, nil
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	var errors []string

	for _, record := range records {
		transformed, err := TransformRecord(record)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Record %s: %v", record.ID, err))
			continue
		}
		processed = append(processed, transformed)
	}

	if len(errors) > 0 {
		return processed, fmt.Errorf("processing completed with errors: %v", strings.Join(errors, "; "))
	}
	return processed, nil
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec1",
			Value:     100.0,
			Timestamp: time.Now(),
			Tags:      []string{"important", "test"},
		},
		{
			ID:        "rec2",
			Value:     -50.0,
			Timestamp: time.Now(),
			Tags:      []string{"invalid"},
		},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
	}

	fmt.Printf("Successfully processed %d records\n", len(processed))
	for _, rec := range processed {
		fmt.Printf("Record: %+v\n", rec)
	}
}
package data_processor

import (
	"regexp"
	"strings"
)

type Processor struct {
	whitespaceRegex *regexp.Regexp
}

func NewProcessor() *Processor {
	return &Processor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

func (p *Processor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	cleaned := p.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

func (p *Processor) NormalizeCase(input string) string {
	return strings.ToLower(p.CleanInput(input))
}

func (p *Processor) ExtractAlphanumeric(input string) string {
	alnumRegex := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return alnumRegex.ReplaceAllString(input, "")
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

	fmt.Printf("Processed %d records, cleaned %d records\n", recordCount, cleanedCount)
	return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
		if cleaned[i] == "" {
			cleaned[i] = "N/A"
		}
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