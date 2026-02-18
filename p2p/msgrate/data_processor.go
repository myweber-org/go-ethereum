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
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func ProcessCSVFile(inputPath string, outputPath string) error {
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

	csvReader := csv.NewReader(inputFile)
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	processedHeader := append(header, "Processed", "Status")
	if err := csvWriter.Write(processedHeader); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	recordCount := 0
	errorCount := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorCount++
			continue
		}

		processedRecord, err := transformRecord(record)
		if err != nil {
			errorCount++
			processedRecord = append(record, "ERROR", err.Error())
		} else {
			recordCount++
			processedRecord = append(record, "SUCCESS", "valid")
		}

		if err := csvWriter.Write(processedRecord); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	fmt.Printf("Processing complete. Records: %d, Errors: %d\n", recordCount, errorCount)
	return nil
}

func transformRecord(record []string) ([]string, error) {
	if len(record) < 4 {
		return nil, fmt.Errorf("insufficient fields in record")
	}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	name := strings.TrimSpace(record[1])
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	value, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value format: %w", err)
	}

	active, err := strconv.ParseBool(record[3])
	if err != nil {
		return nil, fmt.Errorf("invalid active flag format: %w", err)
	}

	data := DataRecord{
		ID:     id,
		Name:   name,
		Value:  value,
		Active: active,
	}

	if data.Value < 0 {
		data.Value = 0
	}

	if !data.Active {
		data.Name = "INACTIVE_" + data.Name
	}

	return []string{
		strconv.Itoa(data.ID),
		data.Name,
		strconv.FormatFloat(data.Value, 'f', 2, 64),
		strconv.FormatBool(data.Active),
	}, nil
}

func ValidateCSVStructure(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	expectedFields := []string{"ID", "Name", "Value", "Active"}
	if len(header) != len(expectedFields) {
		return fmt.Errorf("invalid field count: expected %d, got %d", len(expectedFields), len(header))
	}

	for i, field := range expectedFields {
		if header[i] != field {
			return fmt.Errorf("field mismatch at position %d: expected %s, got %s", i, field, header[i])
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := ValidateCSVStructure(inputFile); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

	if err := ProcessCSVFile(inputFile, outputFile); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Data processing completed successfully")
}
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	ID      int
	Name    string
	Value   float64
	Active  bool
}

func parseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNum, len(row))
		}

		record, err := parseRow(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, errors.New("no valid records found in file")
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (Record, error) {
	var record Record

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return Record{}, fmt.Errorf("empty name at line %d", lineNum)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return Record{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}
	record.Value = value

	active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
	if err != nil {
		return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
	}
	record.Active = active

	return record, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid record ID: %d (must be positive)", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value not allowed for record ID %d", record.ID)
		}
	}

	return nil
}

func processData(filename string) error {
	records, err := parseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if err := validateRecords(records); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("Successfully processed %d records\n", len(records))
	for _, record := range records {
		fmt.Printf("ID: %d, Name: %s, Value: %.2f, Active: %t\n",
			record.ID, record.Name, record.Value, record.Active)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: data_processor <csv_filename>")
		os.Exit(1)
	}

	if err := processData(os.Args[1]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}