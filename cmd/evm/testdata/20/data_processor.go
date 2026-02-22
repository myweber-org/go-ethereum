
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
	validCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		recordCount++
		if dp.isValidRecord(record) {
			cleanedRecord := dp.cleanRecord(record)
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			validCount++
		}
	}

	fmt.Printf("Processed %d records, %d valid records written to %s\n", 
		recordCount, validCount, dp.OutputPath)
	return nil
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) == "" {
			return false
		}
	}
	return true
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
	}
	return cleaned
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
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID    string
    Name  string
    Email string
    Valid bool
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

        if len(row) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(row[0]),
            Name:  strings.TrimSpace(row[1]),
            Email: strings.TrimSpace(row[2]),
            Valid: validateEmail(strings.TrimSpace(row[2])),
        }

        if record.ID != "" && record.Name != "" {
            records = append(records, record)
        }
    }

    return records, nil
}

func validateEmail(email string) bool {
    if email == "" {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func GenerateReport(records []DataRecord) {
    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid email addresses: %d\n", validCount)
    fmt.Printf("Invalid email addresses: %d\n", len(records)-validCount)

    if len(records) > 0 {
        fmt.Println("\nSample records:")
        for i := 0; i < min(3, len(records)); i++ {
            status := "INVALID"
            if records[i].Valid {
                status = "VALID"
            }
            fmt.Printf("  %s: %s (%s)\n", records[i].ID, records[i].Name, status)
        }
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file_path>")
        os.Exit(1)
    }

    filePath := os.Args[1]
    records, err := ProcessCSVFile(filePath)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    GenerateReport(records)
}