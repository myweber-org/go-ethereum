
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
	records := []DataRecord{}
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) != 3 {
			return nil, errors.New("invalid column count at line " + strconv.Itoa(lineNumber))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, errors.New("invalid ID format at line " + strconv.Itoa(lineNumber))
		}

		name := row[1]
		if name == "" {
			return nil, errors.New("empty name at line " + strconv.Itoa(lineNumber))
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, errors.New("invalid value format at line " + strconv.Itoa(lineNumber))
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
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid ID: " + strconv.Itoa(record.ID))
		}
		if seenIDs[record.ID] {
			return errors.New("duplicate ID: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true
	}
	return nil
}
package data

func FilterAndTransform(nums []int, predicate func(int) bool, transform func(int) int) []int {
    var result []int
    for _, n := range nums {
        if predicate(n) {
            result = append(result, transform(n))
        }
    }
    return result
}package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNum := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if lineNum == 0 {
			lineNum++
			continue
		}

		record, err := validateAndCreateRecord(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
		lineNum++
	}

	if len(records) == 0 {
		return nil, errors.New("no valid data records found")
	}

	return records, nil
}

func validateAndCreateRecord(row []string, lineNum int) (DataRecord, error) {
	if len(row) != 4 {
		return DataRecord{}, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNum, len(row))
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}

	name := strings.TrimSpace(row[1])
	if name == "" {
		return DataRecord{}, fmt.Errorf("empty name at line %d", lineNum)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}

	valid, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid valid flag at line %d: %w", lineNum, err)
	}

	return DataRecord{
		ID:    id,
		Name:  name,
		Value: value,
		Valid: valid,
	}, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverageValue(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	total := 0.0
	count := 0
	for _, record := range records {
		if record.Valid {
			total += record.Value
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func GenerateSummary(records []DataRecord) string {
	validCount := 0
	invalidCount := 0
	totalValue := 0.0

	for _, record := range records {
		if record.Valid {
			validCount++
			totalValue += record.Value
		} else {
			invalidCount++
		}
	}

	average := 0.0
	if validCount > 0 {
		average = totalValue / float64(validCount)
	}

	return fmt.Sprintf("Total records: %d, Valid: %d, Invalid: %d, Average value: %.2f",
		len(records), validCount, invalidCount, average)
}