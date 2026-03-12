package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Score float64 `json:"score"`
	Valid bool    `json:"valid"`
}

func processCSVData(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.Comment = '#'

	var records []Record
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNumber, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 columns, got %d", lineNumber, len(row))
		}

		age, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid age: %v", lineNumber, err)
		}

		score, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid score: %v", lineNumber, err)
		}

		valid, err := strconv.ParseBool(row[3])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid valid flag: %v", lineNumber, err)
		}

		record := Record{
			Name:  row[0],
			Age:   age,
			Score: score,
			Valid: valid,
		}
		records = append(records, record)
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func filterValidRecords(records []Record) []Record {
	var filtered []Record
	for _, record := range records {
		if record.Valid {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func calculateAverageScore(records []Record) float64 {
	if len(records) == 0 {
		return 0.0
	}

	total := 0.0
	count := 0
	for _, record := range records {
		if record.Valid {
			total += record.Score
			count++
		}
	}

	if count == 0 {
		return 0.0
	}
	return total / float64(count)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVData(filename)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processed %d total records\n", len(records))

	validRecords := filterValidRecords(records)
	fmt.Printf("Found %d valid records\n", len(validRecords))

	averageScore := calculateAverageScore(validRecords)
	fmt.Printf("Average score of valid records: %.2f\n", averageScore)

	jsonOutput, err := convertToJSON(validRecords)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	outputFile := "processed_data.json"
	err = os.WriteFile(outputFile, []byte(jsonOutput), 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results written to %s\n", outputFile)
}