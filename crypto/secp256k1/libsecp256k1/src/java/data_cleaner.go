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
	data := []int{1, 2, 2, 3, 4, 4, 5, 6, 6}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func deduplicateRecords(records [][]string) [][]string {
	seen := make(map[string]bool)
	var unique [][]string
	for _, record := range records {
		key := strings.Join(record, "|")
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("empty input file")
	}

	headers := records[0]
	var validRecords [][]string
	validRecords = append(validRecords, headers)

	emailIndex := -1
	for i, header := range headers {
		if strings.ToLower(header) == "email" {
			emailIndex = i
			break
		}
	}

	for _, record := range records[1:] {
		if emailIndex >= 0 && emailIndex < len(record) {
			if !validateEmail(record[emailIndex]) {
				continue
			}
		}
		validRecords = append(validRecords, record)
	}

	cleanedRecords := deduplicateRecords(validRecords[1:])
	cleanedRecords = append([][]string{headers}, cleanedRecords...)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	return writer.WriteAll(cleanedRecords)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := processCSV(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Data cleaning completed successfully")
}