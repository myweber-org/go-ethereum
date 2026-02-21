
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

type DataRecord struct {
	ID    int
	Name  string
	Value float64
	Valid bool
}

func ReadCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		record, err := parseRecord(line, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
		lineNumber++
	}

	return records, nil
}

func parseRecord(fields []string, lineNumber int) (DataRecord, error) {
	if len(fields) != 4 {
		return DataRecord{}, fmt.Errorf("invalid field count at line %d: expected 4, got %d", lineNumber, len(fields))
	}

	id, err := strconv.Atoi(strings.TrimSpace(fields[0]))
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}

	name := strings.TrimSpace(fields[1])
	if name == "" {
		return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
	if err != nil {
		return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}

	valid := false
	if strings.ToLower(strings.TrimSpace(fields[3])) == "true" {
		valid = true
	} else if strings.ToLower(strings.TrimSpace(fields[3])) != "false" {
		return DataRecord{}, fmt.Errorf("invalid boolean at line %d: %s", lineNumber, fields[3])
	}

	return DataRecord{
		ID:    id,
		Name:  name,
		Value: value,
		Valid: valid,
	}, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var validRecords []DataRecord
	for _, record := range records {
		if record.Valid {
			validRecords = append(validRecords, record)
		}
	}
	return validRecords
}

func CalculateAverage(records []DataRecord) (float64, error) {
	if len(records) == 0 {
		return 0, errors.New("no records to calculate average")
	}

	var sum float64
	count := 0
	for _, record := range records {
		if record.Valid {
			sum += record.Value
			count++
		}
	}

	if count == 0 {
		return 0, errors.New("no valid records to calculate average")
	}

	return sum / float64(count), nil
}

func FindMaxValue(records []DataRecord) (DataRecord, error) {
	if len(records) == 0 {
		return DataRecord{}, errors.New("no records to find maximum")
	}

	var maxRecord DataRecord
	found := false

	for _, record := range records {
		if record.Valid && (!found || record.Value > maxRecord.Value) {
			maxRecord = record
			found = true
		}
	}

	if !found {
		return DataRecord{}, errors.New("no valid records to find maximum")
	}

	return maxRecord, nil
}package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if len(strings.TrimSpace(data.Username)) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: transformedUsername,
		Age:      age,
	}
	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}
	return userData, nil
}