
package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Email     string
	Timestamp time.Time
	Value     float64
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeString(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func ProcessRecord(record DataRecord) (DataRecord, error) {
	if err := ValidateEmail(record.Email); err != nil {
		return DataRecord{}, err
	}

	record.Email = NormalizeString(record.Email)
	record.ID = strings.ToUpper(record.ID)

	if record.Value < 0 {
		record.Value = 0
	}

	return record, nil
}

func FilterRecords(records []DataRecord, minValue float64) []DataRecord {
	var filtered []DataRecord
	for _, record := range records {
		if record.Value >= minValue {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func CalculateAverage(records []DataRecord) float64 {
	if len(records) == 0 {
		return 0
	}

	var sum float64
	for _, record := range records {
		sum += record.Value
	}
	return sum / float64(len(records))
}
package data

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
	return p.allowedPattern.FindString(trimmed)
}

func (p *Processor) Validate(input string) bool {
	return p.allowedPattern.MatchString(input)
}

func (p *Processor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		cleaned := p.CleanInput(input)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}