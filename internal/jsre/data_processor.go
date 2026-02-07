package data_processor

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type RecordValidator func([]string) error

func ProcessCSVFile(filePath string, validator RecordValidator) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	normalizedHeaders := normalizeHeaders(headers)
	var results []map[string]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if validator != nil {
			if err := validator(record); err != nil {
				continue
			}
		}

		recordMap := make(map[string]string)
		for i, header := range normalizedHeaders {
			if i < len(record) {
				recordMap[header] = strings.TrimSpace(record[i])
			}
		}
		results = append(results, recordMap)
	}

	return results, nil
}

func normalizeHeaders(headers []string) []string {
	normalized := make([]string, len(headers))
	for i, header := range headers {
		normalized[i] = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(header), " ", "_"))
	}
	return normalized
}

func ValidateRecordLength(expected int) RecordValidator {
	return func(record []string) error {
		if len(record) != expected {
			return errors.New("record length mismatch")
		}
		return nil
	}
}