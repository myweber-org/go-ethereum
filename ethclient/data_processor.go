
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) NormalizeEmail(email string) (string, bool) {
	cleaned := dp.CleanString(email)
	normalized := strings.ToLower(cleaned)
	return normalized, dp.ValidateEmail(normalized)
}

func (dp *DataProcessor) ProcessInputList(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		cleaned := dp.CleanString(input)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}
package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []DataRecord
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		if len(line) != 3 {
			return nil, errors.New("invalid column count at line " + strconv.Itoa(lineNumber))
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, errors.New("invalid ID at line " + strconv.Itoa(lineNumber))
		}

		name := line[1]
		if name == "" {
			return nil, errors.New("empty name at line " + strconv.Itoa(lineNumber))
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, errors.New("invalid value at line " + strconv.Itoa(lineNumber))
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)
	for _, record := range records {
		if record.ID <= 0 {
			return errors.New("invalid ID: " + strconv.Itoa(record.ID))
		}
		if seenIDs[record.ID] {
			return errors.New("duplicate ID: " + strconv.Itoa(record.ID))
		}
		seenIDs[record.ID] = true
	}
	return nil
}