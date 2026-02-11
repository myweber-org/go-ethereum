
package main

import "fmt"

func calculateAverage(numbers []int) float64 {
    if len(numbers) == 0 {
        return 0
    }
    
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    
    return float64(sum) / float64(len(numbers))
}

func main() {
    data := []int{10, 20, 30, 40, 50}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
}
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filePath string) ([]DataRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]DataRecord, 0)

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
			return nil, errors.New("invalid CSV format: insufficient columns")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, errors.New("invalid ID format")
		}

		name := row[1]
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.New("invalid value format")
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) error {
	if len(records) == 0 {
		return errors.New("no records to validate")
	}

	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("ID must be positive integer")
		}
		if seenIDs[record.ID] {
			return errors.New("duplicate ID found: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return errors.New("value cannot be negative for ID: " + strconv.Itoa(record.ID))
		}
	}

	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}