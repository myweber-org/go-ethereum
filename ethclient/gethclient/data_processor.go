
package main

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
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
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

func ProcessCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNumber := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber+1, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber+1, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber+1, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber+1)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber+1, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
		lineNumber++
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStats(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	avg, max, count := CalculateStats(records)
	fmt.Printf("Processed %d records\n", count)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)

	for i, record := range records {
		if i < 3 {
			fmt.Printf("Sample record %d: ID=%d, Name=%s, Value=%.2f\n", i+1, record.ID, record.Name, record.Value)
		}
	}
}