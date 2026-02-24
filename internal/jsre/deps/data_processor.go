
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
	ID        int
	Name      string
	Value     float64
	Timestamp string
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []DataRecord
	lineNumber := 0

	for {
		lineNumber++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
		}

		if len(row) != 4 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 4, got %d", lineNumber, len(row))
		}

		record, err := parseRow(row, lineNumber)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, errors.New("no valid records found in file")
	}

	return records, nil
}

func parseRow(row []string, lineNumber int) (DataRecord, error) {
	var record DataRecord

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return record, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
	}
	record.ID = id

	record.Name = strings.TrimSpace(row[1])
	if record.Name == "" {
		return record, fmt.Errorf("empty name at line %d", lineNumber)
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
	if err != nil {
		return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
	}
	record.Value = value

	record.Timestamp = strings.TrimSpace(row[3])
	if record.Timestamp == "" {
		return record, fmt.Errorf("empty timestamp at line %d", lineNumber)
	}

	return record, nil
}

func ValidateRecords(records []DataRecord) error {
	seenIDs := make(map[int]bool)

	for _, record := range records {
		if record.ID <= 0 {
			return fmt.Errorf("invalid ID %d: must be positive", record.ID)
		}

		if seenIDs[record.ID] {
			return fmt.Errorf("duplicate ID %d found", record.ID)
		}
		seenIDs[record.ID] = true

		if record.Value < 0 {
			return fmt.Errorf("negative value %f for record ID %d", record.Value, record.ID)
		}
	}

	return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, float64) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var min, max float64

	for i, record := range records {
		sum += record.Value

		if i == 0 {
			min = record.Value
			max = record.Value
		} else {
			if record.Value < min {
				min = record.Value
			}
			if record.Value > max {
				max = record.Value
			}
		}
	}

	average := sum / float64(len(records))
	return average, min, max
}package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Score  float64 `json:"score"`
	Active bool    `json:"active"`
}

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var records []Record
	lineNumber := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", lineNumber, err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		if len(line) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 columns, got %d", lineNumber, len(line))
		}

		age, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid age: %v", lineNumber, err)
		}

		score, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid score: %v", lineNumber, err)
		}

		active, err := strconv.ParseBool(line[3])
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid active flag: %v", lineNumber, err)
		}

		record := Record{
			Name:   line[0],
			Age:    age,
			Score:  score,
			Active: active,
		}

		records = append(records, record)
		lineNumber++
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func calculateStats(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var totalScore float64
	var totalAge int
	activeCount := 0

	for _, record := range records {
		totalScore += record.Score
		totalAge += record.Age
		if record.Active {
			activeCount++
		}
	}

	averageScore := totalScore / float64(len(records))
	averageAge := float64(totalAge) / float64(len(records))

	return averageScore, averageAge, activeCount
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := processCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	jsonOutput, err := convertToJSON(records)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("JSON Output:")
	fmt.Println(jsonOutput)

	avgScore, avgAge, activeCount := calculateStats(records)
	fmt.Printf("\nStatistics:\n")
	fmt.Printf("Average Score: %.2f\n", avgScore)
	fmt.Printf("Average Age: %.2f\n", avgAge)
	fmt.Printf("Active Records: %d\n", activeCount)
	fmt.Printf("Total Records: %d\n", len(records))
}
package main

import (
    "errors"
    "regexp"
    "strings"
)

type UserProfile struct {
    ID        string
    Email     string
    Username  string
    Age       int
    IsActive  bool
}

func ValidateEmail(email string) error {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, err := regexp.MatchString(pattern, email)
    if err != nil {
        return err
    }
    if !matched {
        return errors.New("invalid email format")
    }
    return nil
}

func NormalizeUsername(username string) string {
    return strings.ToLower(strings.TrimSpace(username))
}

func TransformUserProfile(profile UserProfile) (UserProfile, error) {
    if err := ValidateEmail(profile.Email); err != nil {
        return profile, err
    }

    normalizedUsername := NormalizeUsername(profile.Username)

    transformed := UserProfile{
        ID:        strings.ToUpper(profile.ID),
        Email:     strings.ToLower(profile.Email),
        Username:  normalizedUsername,
        Age:       profile.Age,
        IsActive:  profile.IsActive,
    }

    if transformed.Age < 0 {
        return transformed, errors.New("age cannot be negative")
    }

    return transformed, nil
}

func ProcessUserBatch(profiles []UserProfile) ([]UserProfile, []error) {
    var processed []UserProfile
    var errs []error

    for i, profile := range profiles {
        transformed, err := TransformUserProfile(profile)
        if err != nil {
            errs = append(errs, errors.New("profile index "+string(rune(i))+": "+err.Error()))
            continue
        }
        processed = append(processed, transformed)
    }

    return processed, errs
}package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

func processCSVFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) < 3 {
			continue
		}

		value, err1 := strconv.ParseFloat(row[1], 64)
		count, err2 := strconv.Atoi(row[2])

		if err1 != nil || err2 != nil {
			continue
		}

		record := Record{
			Name:  row[0],
			Value: value,
			Count: count,
		}
		records = append(records, record)
	}

	return records, nil
}

func convertToJSON(records []Record) (string, error) {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	inputFile := os.Args[1]
	records, err := processCSVFile(inputFile)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	jsonOutput, err := convertToJSON(records)
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		return
	}

	fmt.Println("Processed Data:")
	fmt.Println(jsonOutput)
	fmt.Printf("Total records processed: %d\n", len(records))
}