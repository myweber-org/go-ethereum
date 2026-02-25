package data_processor

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	var builder strings.Builder
	for _, r := range input {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func TransformToSlug(text string) string {
	text = strings.ToLower(text)
	text = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(text, "-")
	text = strings.Trim(text, "-")
	return text
}

func ExtractDomain(email string) (string, error) {
	if !ValidateEmail(email) {
		return "", errors.New("invalid email format")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", errors.New("malformed email address")
	}
	return parts[1], nil
}package main

import (
	"encoding/csv"
	"errors"
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

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []DataRecord
	for i, row := range records {
		if len(row) < 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(row[1])
		if name == "" {
			continue
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
		if err != nil {
			continue
		}

		valid := strings.ToLower(strings.TrimSpace(row[3])) == "true"

		record := DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
			Valid: valid,
		}
		result = append(result, record)
	}

	return result, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, error) {
	if len(records) == 0 {
		return nil, errors.New("no records to validate")
	}

	var validRecords []DataRecord
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			continue
		}
		if seenIDs[record.ID] {
			continue
		}
		if record.Value < 0 {
			continue
		}

		seenIDs[record.ID] = true
		validRecords = append(validRecords, record)
	}

	return validRecords, nil
}

func WriteProcessedData(records []DataRecord, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Value", "Valid"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, record := range records {
		row := []string{
			strconv.Itoa(record.ID),
			record.Name,
			strconv.FormatFloat(record.Value, 'f', 2, 64),
			strconv.FormatBool(record.Valid),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func ProcessDataPipeline(inputFile, outputFile string) error {
	records, err := ParseCSVFile(inputFile)
	if err != nil {
		return err
	}

	validRecords, err := ValidateRecords(records)
	if err != nil {
		return err
	}

	if len(validRecords) == 0 {
		return errors.New("no valid records after processing")
	}

	return WriteProcessedData(validRecords, outputFile)
}
package main

import (
	"regexp"
	"strings"
)

var (
	invalidCharRegex = regexp.MustCompile(`[<>"'&;]`)
	whitespaceRegex  = regexp.MustCompile(`\s+`)
)

func SanitizeInput(input string) string {
	if input == "" {
		return input
	}

	trimmed := strings.TrimSpace(input)
	sanitized := invalidCharRegex.ReplaceAllString(trimmed, "")
	normalized := whitespaceRegex.ReplaceAllString(sanitized, " ")

	return normalized
}

func ValidateInput(input string, maxLength int) bool {
	if input == "" || len(input) > maxLength {
		return false
	}

	sanitized := SanitizeInput(input)
	return sanitized == input
}package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func ProcessUserInput(rawUsername string, rawEmail string, rawAge int) (UserData, error) {
	userData := UserData{
		Username: rawUsername,
		Email:    rawEmail,
		Age:      rawAge,
	}

	userData = TransformUsername(userData)

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}