
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	filePath string
	headers  []string
}

func NewDataProcessor(filePath string) *DataProcessor {
	return &DataProcessor{
		filePath: filePath,
	}
}

func (dp *DataProcessor) ValidateHeaders(expected []string) error {
	file, err := os.Open(dp.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if len(headers) != len(expected) {
		return fmt.Errorf("header count mismatch: expected %d, got %d", len(expected), len(headers))
	}

	for i, header := range headers {
		if strings.TrimSpace(header) != expected[i] {
			return fmt.Errorf("header mismatch at position %d: expected '%s', got '%s'", i, expected[i], header)
		}
	}

	dp.headers = headers
	return nil
}

func (dp *DataProcessor) ProcessRows(handler func(row []string) error) error {
	file, err := os.Open(dp.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = len(dp.headers)

	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to skip headers: %w", err)
	}

	rowNumber := 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading row %d: %w", rowNumber, err)
		}

		if err := handler(row); err != nil {
			return fmt.Errorf("error processing row %d: %w", rowNumber, err)
		}
		rowNumber++
	}

	return nil
}

func (dp *DataProcessor) CountRows() (int, error) {
	count := 0
	err := dp.ProcessRows(func(row []string) error {
		count++
		return nil
	})
	return count, err
}