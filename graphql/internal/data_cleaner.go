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

	headers = append(headers, "Validated", "Grade")
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
			fmt.Printf("warning: line %d: %v\n", lineNum, err)
			continue
		}

		record, validationErr := validateRecord(row)
		validated := "PASS"
		if validationErr != nil {
			validated = "FAIL: " + validationErr.Error()
		}

		grade := calculateGrade(record.Score)
		outputRow := []string{
			strconv.Itoa(record.ID),
			strings.TrimSpace(record.Name),
			strings.ToLower(strings.TrimSpace(record.Email)),
			strconv.FormatBool(record.Active),
			fmt.Sprintf("%.2f", record.Score),
			validated,
			grade,
		}

		if err := writer.Write(outputRow); err != nil {
			fmt.Printf("warning: failed to write line %d: %v\n", lineNum, err)
		}
	}

	return nil
}

func validateRecord(row []string) (Record, error) {
	if len(row) < 5 {
		return Record{}, fmt.Errorf("insufficient columns")
	}

	var rec Record
	var err error

	rec.ID, err = strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID: %w", err)
	}

	rec.Name = row[1]
	if strings.TrimSpace(rec.Name) == "" {
		return Record{}, fmt.Errorf("empty name")
	}

	rec.Email = row[2]
	if !strings.Contains(rec.Email, "@") {
		return Record{}, fmt.Errorf("invalid email format")
	}

	rec.Active, err = strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag: %w", err)
	}

	rec.Score, err = strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid score: %w", err)
	}

	if rec.Score < 0 || rec.Score > 100 {
		return Record{}, fmt.Errorf("score out of range (0-100)")
	}

	return rec, nil
}

func calculateGrade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
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

	fmt.Printf("Data cleaning completed. Output written to %s\n", outputFile)
}