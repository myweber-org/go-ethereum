
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

func ParseCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]Record, 0)

    headerSkipped := false
    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        if !headerSkipped {
            headerSkipped = true
            continue
        }

        if len(row) != 4 {
            return nil, errors.New("invalid row format")
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID format: %w", err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value format: %w", err)
        }

        active, err := strconv.ParseBool(row[3])
        if err != nil {
            return nil, fmt.Errorf("invalid active flag format: %w", err)
        }

        record := Record{
            ID:     id,
            Name:   row[1],
            Value:  value,
            Active: active,
        }
        records = append(records, record)
    }

    return records, nil
}

func ValidateRecords(records []Record) error {
    if len(records) == 0 {
        return errors.New("no records found")
    }

    seenIDs := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid ID %d: must be positive", record.ID)
        }

        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        seenIDs[record.ID] = true

        if record.Name == "" {
            return fmt.Errorf("empty name for record ID %d", record.ID)
        }

        if record.Value < 0 {
            return fmt.Errorf("negative value for record ID %d", record.ID)
        }
    }

    return nil
}

func CalculateStatistics(records []Record) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var max float64
    activeCount := 0

    for i, record := range records {
        sum += record.Value
        if i == 0 || record.Value > max {
            max = record.Value
        }
        if record.Active {
            activeCount++
        }
    }

    average := sum / float64(len(records))
    return average, max, activeCount
}package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type Record struct {
	ID    int
	Name  string
	Value float64
}

func ReadCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := make([]Record, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil && err != io.EOF {
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
			return nil, errors.New("invalid csv format")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []Record) error {
	if len(records) == 0 {
		return errors.New("no records found")
	}

	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid id: " + strconv.Itoa(record.ID))
		}

		if seenIDs[record.ID] {
			return errors.New("duplicate id: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true

		if record.Name == "" {
			return errors.New("empty name for id: " + strconv.Itoa(record.ID))
		}

		if record.Value < 0 {
			return errors.New("negative value for id: " + strconv.Itoa(record.ID))
		}
	}

	return nil
}

func CalculateTotalValue(records []Record) float64 {
	total := 0.0
	for _, record := range records {
		total += record.Value
	}
	return total
}