
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

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) < 4 {
			return nil, fmt.Errorf("invalid column count at line %d", lineNum)
		}

		record, err := parseRow(row, lineNum)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func parseRow(row []string, lineNum int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
	}
	record.ID = id

	name := strings.TrimSpace(row[1])
	if name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNum)
	}
	record.Name = name

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
	}
	record.Value = value

	validStr := strings.ToLower(strings.TrimSpace(row[3]))
	if validStr != "true" && validStr != "false" {
		return record, fmt.Errorf("invalid boolean at line %d", lineNum)
	}
	record.Valid = (validStr == "true")

	return record, nil
}

func ValidateRecords(records []DataRecord) ([]DataRecord, error) {
	if len(records) == 0 {
		return nil, errors.New("no records to validate")
	}

	validRecords := []DataRecord{}
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if seenIDs[record.ID] {
			return nil, fmt.Errorf("duplicate ID found: %d", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			continue
		}

		validRecords = append(validRecords, record)
	}

	return validRecords, nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var count int
	var maxValue float64

	for _, record := range records {
		if record.Valid {
			sum += record.Value
			count++
			if record.Value > maxValue {
				maxValue = record.Value
			}
		}
	}

	if count == 0 {
		return 0, 0, 0
	}

	average := sum / float64(count)
	return average, maxValue, count
}

func ProcessDataFile(filename string) error {
	records, err := ParseCSVFile(filename)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	validRecords, err := ValidateRecords(records)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	avg, max, count := CalculateStatistics(validRecords)
	fmt.Printf("Processed %d valid records\n", len(validRecords))
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
	fmt.Printf("Valid records used for stats: %d\n", count)

	return nil
}
package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func sanitizeUsername(username string) string {
	username = strings.TrimSpace(username)
	username = strings.ToLower(username)
	return username
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if !validateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = sanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 120 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	jsonData := []byte(`{"email":"test@example.com","username":"  JohnDoe  ","age":25}`)
	processedData, err := processUserData(jsonData)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processedData)
}