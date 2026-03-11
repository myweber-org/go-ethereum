package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Bio      string
}

func SanitizeInput(input string) string {
	// Remove leading/trailing whitespace
	trimmed := strings.TrimSpace(input)
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(trimmed, " ")
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

func ProcessUserData(data UserData) (UserData, error) {
	sanitizedData := UserData{
		Username: SanitizeInput(data.Username),
		Email:    SanitizeInput(data.Email),
		Bio:      SanitizeInput(data.Bio),
	}

	if !ValidateUsername(sanitizedData.Username) {
		return UserData{}, &ValidationError{Field: "username", Message: "invalid username format"}
	}

	if !ValidateEmail(sanitizedData.Email) {
		return UserData{}, &ValidationError{Field: "email", Message: "invalid email format"}
	}

	if len(sanitizedData.Bio) > 500 {
		sanitizedData.Bio = sanitizedData.Bio[:500]
	}

	return sanitizedData, nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
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
	records := make([]DataRecord, 0)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			return nil, errors.New("invalid CSV format")
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}

		name := row[1]

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, err
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
			return errors.New("invalid ID: must be positive integer")
		}

		if seenIDs[record.ID] {
			return errors.New("duplicate ID found")
		}
		seenIDs[record.ID] = true

		if record.Name == "" {
			return errors.New("name cannot be empty")
		}

		if record.Value < 0 {
			return errors.New("value cannot be negative")
		}
	}

	return nil
}

func CalculateTotalValue(records []DataRecord) float64 {
	total := 0.0
	for _, record := range records {
		total += record.Value
	}
	return total
}