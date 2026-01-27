
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
)

type DataRecord struct {
	ID   string
	Name string
	Age  string
}

func readCSV(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []DataRecord

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(row) >= 3 {
			records = append(records, DataRecord{
				ID:   row[0],
				Name: row[1],
				Age:  row[2],
			})
		}
	}
	return records, nil
}

func removeDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := record.ID + "|" + record.Name
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func sortByAge(records []DataRecord) []DataRecord {
	sorted := make([]DataRecord, len(records))
	copy(sorted, records)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Age < sorted[j].Age
	})
	return sorted
}

func writeCSV(filename string, records []DataRecord) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		row := []string{record.ID, record.Name, record.Age}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	inputFile := "input.csv"
	outputFile := "cleaned_data.csv"

	records, err := readCSV(inputFile)
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	fmt.Printf("Original records: %d\n", len(records))

	uniqueRecords := removeDuplicates(records)
	fmt.Printf("After duplicate removal: %d\n", len(uniqueRecords))

	sortedRecords := sortByAge(uniqueRecords)

	if err := writeCSV(outputFile, sortedRecords); err != nil {
		fmt.Printf("Error writing CSV: %v\n", err)
		return
	}

	fmt.Printf("Cleaned data written to %s\n", outputFile)
}