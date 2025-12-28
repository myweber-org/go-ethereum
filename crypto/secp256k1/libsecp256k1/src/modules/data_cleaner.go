
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func sanitizeCSVRow(row []string) []string {
	sanitized := make([]string, len(row))
	for i, field := range row {
		trimmed := strings.TrimSpace(field)
		normalized := strings.ToLower(trimmed)
		sanitized[i] = normalized
	}
	return sanitized
}

func processCSV(reader io.Reader, writer io.Writer) error {
	csvReader := csv.NewReader(reader)
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		cleanRecord := sanitizeCSVRow(record)
		if err := csvWriter.Write(cleanRecord); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}
	return nil
}