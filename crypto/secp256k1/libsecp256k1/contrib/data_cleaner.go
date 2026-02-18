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
	ID    int
	Name  string
	Email string
	Age   int
}

func cleanEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateAge(age int) bool {
	return age >= 0 && age <= 120
}

func parseCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lineNum++

		if lineNum == 1 {
			continue
		}

		if len(line) < 4 {
			continue
		}

		id, err1 := strconv.Atoi(strings.TrimSpace(line[0]))
		age, err2 := strconv.Atoi(strings.TrimSpace(line[3]))

		if err1 != nil || err2 != nil {
			continue
		}

		if !validateAge(age) {
			continue
		}

		record := Record{
			ID:    id,
			Name:  strings.TrimSpace(line[1]),
			Email: cleanEmail(line[2]),
			Age:   age,
		}
		records = append(records, record)
	}

	return records, nil
}

func removeDuplicates(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, record := range records {
		key := record.Email
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func generateReport(records []Record) {
	fmt.Printf("Total valid records: %d\n", len(records))
	fmt.Println("ID | Name | Email | Age")
	fmt.Println("---|------|-------|----")
	for _, record := range records {
		fmt.Printf("%d | %s | %s | %d\n", record.ID, record.Name, record.Email, record.Age)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run data_cleaner.go <csv_file>")
		return
	}

	records, err := parseCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	uniqueRecords := removeDuplicates(records)
	generateReport(uniqueRecords)
}