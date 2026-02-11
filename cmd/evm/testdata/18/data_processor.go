
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor(pattern string) (*DataProcessor, error) {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &DataProcessor{allowedPattern: compiled}, nil
}

func (dp *DataProcessor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.allowedPattern.FindString(trimmed)
}

func (dp *DataProcessor) Validate(input string) bool {
	return dp.allowedPattern.MatchString(input)
}

func (dp *DataProcessor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		cleaned := dp.CleanInput(input)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}
package main

import (
	"errors"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Category  string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.After(time.Now()) {
		return errors.New("timestamp cannot be in the future")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Category = strings.ToUpper(strings.TrimSpace(record.Category))
	transformed.Value = roundToTwoDecimals(record.Value)
	return transformed
}

func roundToTwoDecimals(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, err
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}