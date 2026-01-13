package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func CleanCSVRow(row []string) []string {
	cleaned := make([]string, len(row))
	for i, field := range row {
		trimmed := strings.TrimSpace(field)
		normalized := strings.ToLower(trimmed)
		cleaned[i] = normalized
	}
	return cleaned
}

func ProcessCSV(reader io.Reader, writer io.Writer) error {
	csvReader := csv.NewReader(reader)
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := CleanCSVRow(record)
		if err := csvWriter.Write(cleanedRecord); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	sampleInput := "Name,Age,Email\nJohn Doe,25,JOHN@example.com\n Jane Smith,30,jane@test.org "
	reader := strings.NewReader(sampleInput)
	var output strings.Builder

	err := ProcessCSV(reader, &output)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		return
	}

	fmt.Println("Cleaned CSV output:")
	fmt.Println(output.String())
}