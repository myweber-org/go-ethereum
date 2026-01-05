
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func processCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func validateRecords(records []Record) []Record {
	var validRecords []Record
	for _, r := range records {
		if r.ID > 0 && r.Name != "" && r.Value >= 0 {
			validRecords = append(validRecords, r)
		}
	}
	return validRecords
}

func calculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	records, err := processCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	validRecords := validateRecords(records)
	total := calculateTotal(validRecords)

	fmt.Printf("Total records: %d\n", len(records))
	fmt.Printf("Valid records: %d\n", len(validRecords))
	fmt.Printf("Total value: %.2f\n", total)
}