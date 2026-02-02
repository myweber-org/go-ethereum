
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
    Delimiter  rune
}

func NewDataProcessor(input, output string) *DataProcessor {
    return &DataProcessor{
        InputPath:  input,
        OutputPath: output,
        Delimiter:  ',',
    }
}

func (dp *DataProcessor) ValidateRow(row []string) bool {
    if len(row) == 0 {
        return false
    }
    for _, field := range row {
        if strings.TrimSpace(field) == "" {
            return false
        }
    }
    return true
}

func (dp *DataProcessor) CleanField(field string) string {
    cleaned := strings.TrimSpace(field)
    cleaned = strings.ToLower(cleaned)
    return cleaned
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
    reader.Comma = dp.Delimiter

    writer := csv.NewWriter(outputFile)
    writer.Comma = dp.Delimiter
    defer writer.Flush()

    lineCount := 0
    processedCount := 0

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV: %w", err)
        }

        lineCount++

        if !dp.ValidateRow(record) {
            continue
        }

        cleanedRecord := make([]string, len(record))
        for i, field := range record {
            cleanedRecord[i] = dp.CleanField(field)
        }

        if err := writer.Write(cleanedRecord); err != nil {
            return fmt.Errorf("error writing CSV: %w", err)
        }

        processedCount++
    }

    fmt.Printf("Processing complete. Read %d lines, wrote %d valid records.\n", lineCount, processedCount)
    return nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.csv>")
        os.Exit(1)
    }

    processor := NewDataProcessor(os.Args[1], os.Args[2])
    if err := processor.Process(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Category  string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformCategory(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToUpper(trimmed)
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0.0
	}
	
	var sum float64
	validCount := 0
	
	for _, record := range records {
		if err := ValidateRecord(record); err == nil {
			sum += record.Value
			validCount++
		}
	}
	
	if validCount == 0 {
		return 0.0
	}
	return sum / float64(validCount)
}

func FilterByCategory(records []DataRecord, category string) []DataRecord {
	var filtered []DataRecord
	targetCategory := TransformCategory(category)
	
	for _, record := range records {
		if TransformCategory(record.Category) == targetCategory {
			filtered = append(filtered, record)
		}
	}
	return filtered
}