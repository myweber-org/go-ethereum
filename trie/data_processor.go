
package data_processor

import (
	"regexp"
	"strings"
)

type Processor struct {
	allowedPattern *regexp.Regexp
}

func NewProcessor(allowedPattern string) (*Processor, error) {
	compiled, err := regexp.Compile(allowedPattern)
	if err != nil {
		return nil, err
	}
	return &Processor{allowedPattern: compiled}, nil
}

func (p *Processor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	if p.allowedPattern == nil {
		return trimmed
	}
	return p.allowedPattern.FindString(trimmed)
}

func (p *Processor) ValidateInput(input string) bool {
	if input == "" {
		return false
	}
	if p.allowedPattern != nil && !p.allowedPattern.MatchString(input) {
		return false
	}
	return true
}

func (p *Processor) Process(input string) (string, bool) {
	cleaned := p.CleanInput(input)
	valid := p.ValidateInput(cleaned)
	return cleaned, valid
}
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

func parseCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record

	for i := 0; ; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if i == 0 {
			continue
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid row length at line %d", i+1)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %v", i+1, err)
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %v", i+1, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func validateRecords(records []Record) error {
	seenIDs := make(map[int]bool)
	for _, r := range records {
		if r.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", r.ID)
		}
		if seenIDs[r.ID] {
			return fmt.Errorf("duplicate ID %d found", r.ID)
		}
		seenIDs[r.ID] = true

		if r.Name == "" {
			return fmt.Errorf("empty name for ID %d", r.ID)
		}

		if r.Value < 0 {
			return fmt.Errorf("negative value %f for ID %d", r.Value, r.ID)
		}
	}
	return nil
}

func processData(filename string) ([]Record, error) {
	records, err := parseCSV(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %v", err)
	}

	if err := validateRecords(records); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	return records, nil
}