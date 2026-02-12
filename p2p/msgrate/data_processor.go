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

func ReadCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) != 3 {
			continue
		}

		id, err1 := strconv.Atoi(row[0])
		value, err2 := strconv.ParseFloat(row[2], 64)

		if err1 != nil || err2 != nil {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  row[1],
			Value: value,
		})
	}

	return records, nil
}

func CalculateTotal(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}

func FilterByThreshold(records []Record, threshold float64) []Record {
	var filtered []Record
	for _, r := range records {
		if r.Value >= threshold {
			filtered = append(filtered, r)
		}
	}
	return filtered
}