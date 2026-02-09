
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataRecord struct {
	ID      string
	Name    string
	Email   string
	Active  string
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if lineNumber == 1 {
			continue
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
		}

		record := DataRecord{
			ID:     strings.TrimSpace(row[0]),
			Name:   strings.TrimSpace(row[1]),
			Email:  strings.TrimSpace(row[2]),
			Active: strings.TrimSpace(row[3]),
		}

		if record.ID == "" || record.Name == "" || record.Email == "" {
			return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
		}

		if record.Active != "true" && record.Active != "false" {
			return nil, fmt.Errorf("invalid active status at line %d: %s", lineNumber, record.Active)
		}

		records = append(records, record)
	}

	return records, nil
}

func ValidateEmailFormat(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterActiveRecords(records []DataRecord) []DataRecord {
	var active []DataRecord
	for _, record := range records {
		if record.Active == "true" {
			active = append(active, record)
		}
	}
	return active
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total records processed: %d\n", len(records))
	active := FilterActiveRecords(records)
	fmt.Printf("Active records: %d\n", len(active))
	fmt.Printf("Inactive records: %d\n", len(records)-len(active))

	emailValid := 0
	for _, record := range records {
		if ValidateEmailFormat(record.Email) {
			emailValid++
		}
	}
	fmt.Printf("Valid email addresses: %d\n", emailValid)
}package main

import (
	"encoding/csv"
	"errors"
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

func ParseCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNum++
		if lineNum == 1 {
			continue
		}

		record, parseErr := parseLine(line, lineNum)
		if parseErr != nil {
			return nil, parseErr
		}
		records = append(records, record)
	}

	return records, nil
}

func parseLine(fields []string, lineNum int) (Record, error) {
	if len(fields) != 4 {
		return Record{}, errors.New("invalid field count at line " + strconv.Itoa(lineNum))
	}

	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return Record{}, errors.New("invalid ID at line " + strconv.Itoa(lineNum))
	}

	name := strings.TrimSpace(fields[1])
	if name == "" {
		return Record{}, errors.New("empty name at line " + strconv.Itoa(lineNum))
	}

	value, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Record{}, errors.New("invalid value at line " + strconv.Itoa(lineNum))
	}

	active, err := strconv.ParseBool(fields[3])
	if err != nil {
		return Record{}, errors.New("invalid active flag at line " + strconv.Itoa(lineNum))
	}

	return Record{
		ID:     id,
		Name:   name,
		Value:  value,
		Active: active,
	}, nil
}

func FilterActiveRecords(records []Record) []Record {
	var active []Record
	for _, r := range records {
		if r.Active {
			active = append(active, r)
		}
	}
	return active
}

func CalculateTotalValue(records []Record) float64 {
	var total float64
	for _, r := range records {
		total += r.Value
	}
	return total
}