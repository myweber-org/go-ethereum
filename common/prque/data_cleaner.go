package datautils

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

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(line) != 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(line[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(line[1])
		email := strings.TrimSpace(line[2])
		score, err := strconv.ParseFloat(strings.TrimSpace(line[3]), 64)
		if err != nil {
			continue
		}

		if !validateEmail(email) {
			continue
		}

		record := DataRecord{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		}
		records = append(records, record)
	}

	return records, nil
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func calculateAverageScore(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	total := 0.0
	for _, record := range records {
		total += record.Score
	}
	return total / float64(len(records))
}

func filterHighScorers(records []DataRecord, threshold float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Score >= threshold {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <csv_file>")
		return
	}

	records, err := parseCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		return
	}

	fmt.Printf("Parsed %d valid records\n", len(records))
	fmt.Printf("Average score: %.2f\n", calculateAverageScore(records))

	highScorers := filterHighScorers(records, 80.0)
	fmt.Printf("High scorers (>=80): %d\n", len(highScorers))

	for _, record := range highScorers {
		fmt.Printf("ID: %d, Name: %s, Score: %.1f\n", record.ID, record.Name, record.Score)
	}
}