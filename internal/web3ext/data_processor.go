package main

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

func LoadCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	// Skip header
	if _, err := reader.Read(); err != nil {
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

		if len(row) != 3 {
			return nil, errors.New("invalid row length")
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
	seen := make(map[int]bool)
	for _, r := range records {
		if r.ID <= 0 {
			return errors.New("invalid ID: must be positive")
		}
		if r.Name == "" {
			return errors.New("empty name")
		}
		if r.Value < 0 {
			return errors.New("negative value not allowed")
		}
		if seen[r.ID] {
			return errors.New("duplicate ID found")
		}
		seen[r.ID] = true
	}
	return nil
}