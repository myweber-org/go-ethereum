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

func cleanCSVData(inputPath string, outputPath string) error {
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

	lineNumber := 1
	for {
		lineNumber++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("line %d: read error: %v\n", lineNumber, err)
			continue
		}

		cleanedRecord, isValid := validateAndCleanRecord(record)
		cleanedRecord = append(cleanedRecord, strconv.FormatBool(isValid))
		
		if err := writer.Write(cleanedRecord); err != nil {
			fmt.Printf("line %d: write error: %v\n", lineNumber, err)
		}
	}

	return nil
}

func validateAndCleanRecord(record []string) ([]string, bool) {
	if len(record) < 4 {
		return record, false
	}

	cleaned := make([]string, len(record))
	
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
	cleaned[1] = name

	email := strings.ToLower(strings.TrimSpace(record[2]))
	if !strings.Contains(email, "@") {
		email = "invalid@example.com"
	}
	cleaned[2] = email

	score, err := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
	if err != nil || score < 0 || score > 100 {
		cleaned[3] = "0.0"
	} else {
		cleaned[3] = fmt.Sprintf("%.2f", score)
	}

	isValid := id > 0 && name != "Unknown" && strings.Contains(email, "@") && score >= 0 && score <= 100
	
	return cleaned, isValid
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