
package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values.
// It takes a slice of data and a window size as input.
// Returns a slice containing the moving averages or an empty slice if window size is invalid.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return []float64{}
	}

	result := make([]float64, len(data)-windowSize+1)
	for i := 0; i <= len(data)-windowSize; i++ {
		sum := 0.0
		for j := i; j < i+windowSize; j++ {
			sum += data[j]
		}
		result[i] = sum / float64(windowSize)
	}
	return result
}

func main() {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := CalculateMovingAverage(data, window)
	fmt.Printf("Moving averages with window size %d: %v\n", window, averages)
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
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("record value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord, multiplier float64) DataRecord {
	if multiplier <= 0 {
		multiplier = 1.0
	}
	return DataRecord{
		ID:        strings.ToUpper(record.ID),
		Value:     record.Value * multiplier,
		Timestamp: record.Timestamp.UTC(),
		Tags:      append([]string{"processed"}, record.Tags...),
	}
}

func ProcessRecords(records []DataRecord, multiplier float64) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, err
		}
		processed = append(processed, TransformRecord(record, multiplier))
	}
	return processed, nil
}
package main

import (
	"errors"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func ValidateUserData(data UserData) error {
	if data.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if strings.TrimSpace(data.Name) == "" {
		return errors.New("name cannot be empty")
	}

	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}

	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}

	return nil
}

func TransformUserData(data UserData) UserData {
	return UserData{
		ID:    data.ID,
		Name:  strings.ToUpper(strings.TrimSpace(data.Name)),
		Email: strings.ToLower(strings.TrimSpace(data.Email)),
		Age:   data.Age,
	}
}

func ProcessUserInput(rawData UserData) (UserData, error) {
	if err := ValidateUserData(rawData); err != nil {
		return UserData{}, err
	}

	processedData := TransformUserData(rawData)
	return processedData, nil
}