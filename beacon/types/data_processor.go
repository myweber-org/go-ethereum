
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

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	for i := 0; ; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if i == 0 {
			continue
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid row length at line %d", i+1)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %v", i+1, err)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %v", i+1, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return records, nil
}

func calculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func validateRecords(records []Record) error {
	seen := make(map[int]bool)
	for _, r := range records {
		if r.ID <= 0 {
			return fmt.Errorf("invalid ID: %d", r.ID)
		}
		if seen[r.ID] {
			return fmt.Errorf("duplicate ID: %d", r.ID)
		}
		seen[r.ID] = true
		if r.Name == "" {
			return fmt.Errorf("empty name for ID: %d", r.ID)
		}
		if r.Value < 0 {
			return fmt.Errorf("negative value for ID: %d", r.ID)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	records, err := parseCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	if err := validateRecords(records); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		os.Exit(1)
	}

	total := calculateTotal(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Total value: %.2f\n", total)
}