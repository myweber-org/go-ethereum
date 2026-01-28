
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Age   int
	Valid bool
}

func cleanCSVData(inputPath string, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	headers = append(headers, "Valid")
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	lineNumber := 1
	for {
		lineNumber++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("warning: line %d: %v\n", lineNumber, err)
			continue
		}

		if len(record) < 4 {
			fmt.Printf("warning: line %d: insufficient columns\n", lineNumber)
			continue
		}

		dataRec := DataRecord{}
		dataRec.ID, err = strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			fmt.Printf("warning: line %d: invalid ID format\n", lineNumber)
			continue
		}

		dataRec.Name = strings.TrimSpace(record[1])
		if dataRec.Name == "" {
			fmt.Printf("warning: line %d: empty name field\n", lineNumber)
			continue
		}

		dataRec.Email = strings.TrimSpace(record[2])
		if !strings.Contains(dataRec.Email, "@") {
			fmt.Printf("warning: line %d: invalid email format\n", lineNumber)
			continue
		}

		dataRec.Age, err = strconv.Atoi(strings.TrimSpace(record[3]))
		if err != nil || dataRec.Age < 0 || dataRec.Age > 120 {
			fmt.Printf("warning: line %d: invalid age value\n", lineNumber)
			continue
		}

		dataRec.Valid = true

		outputRecord := []string{
			strconv.Itoa(dataRec.ID),
			dataRec.Name,
			dataRec.Email,
			strconv.Itoa(dataRec.Age),
			strconv.FormatBool(dataRec.Valid),
		}

		if err := writer.Write(outputRecord); err != nil {
			fmt.Printf("warning: line %d: failed to write record: %v\n", lineNumber, err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSVData(inputFile, outputFile); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("data cleaning completed successfully\n")
}