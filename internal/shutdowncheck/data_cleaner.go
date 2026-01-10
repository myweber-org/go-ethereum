
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
	ID    int
	Name  string
	Email string
	Score float64
}

func cleanCSVData(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
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

	headers = append(headers, "Valid")
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	lineNum := 1
	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading line %d: %w", lineNum, err)
		}

		cleanedRecord, isValid := processRecord(record)
		cleanedRecord = append(cleanedRecord, strconv.FormatBool(isValid))
		
		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("error writing line %d: %w", lineNum, err)
		}
	}

	return nil
}

func processRecord(record []string) ([]string, bool) {
	if len(record) < 4 {
		return padRecord(record, 4), false
	}

	cleaned := make([]string, 4)
	
	id, err := strconv.Atoi(strings.TrimSpace(record[0]))
	if err != nil || id <= 0 {
		cleaned[0] = "0"
	} else {
		cleaned[0] = strconv.Itoa(id)
	}

	name := strings.TrimSpace(record[1])
	if name == "" {
		name = "Unknown"
	}
	cleaned[1] = strings.Title(strings.ToLower(name))

	email := strings.TrimSpace(record[2])
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		email = "invalid@example.com"
	}
	cleaned[2] = strings.ToLower(email)

	score, err := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
	if err != nil || score < 0 || score > 100 {
		cleaned[3] = "0.0"
	} else {
		cleaned[3] = fmt.Sprintf("%.2f", score)
	}

	isValid := cleaned[0] != "0" && cleaned[1] != "Unknown" && 
		!strings.Contains(cleaned[2], "invalid") && cleaned[3] != "0.0"

	return cleaned, isValid
}

func padRecord(record []string, length int) []string {
	padded := make([]string, length)
	for i := 0; i < length; i++ {
		if i < len(record) {
			padded[i] = record[i]
		} else {
			padded[i] = ""
		}
	}
	return padded
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data cleaning completed. Output saved to %s\n", outputFile)
}