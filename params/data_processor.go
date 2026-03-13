
package main

import (
    "encoding/csv"
    "errors"
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

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNum, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNum)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
        }

        if value < 0 {
            return nil, fmt.Errorf("negative value at line %d: %f", lineNum, value)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
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

func ValidateRecords(records []Record) error {
    seenIDs := make(map[int]bool)

    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid record ID: %d", record.ID)
        }

        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        seenIDs[record.ID] = true

        if record.Value > 10000 {
            return fmt.Errorf("value exceeds limit for ID %d: %f", record.ID, record.Value)
        }
    }

    return nil
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

func ProcessCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []Record{}
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		lineNum++
		if lineNum == 1 {
			continue
		}

		if len(line) != 3 {
			return nil, fmt.Errorf("invalid column count on line %d", lineNum)
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID on line %d: %w", lineNum, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("empty name on line %d", lineNum)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value on line %d: %w", lineNum, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func CalculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, r := range records {
		sum += r.Value
		if r.Value > max {
			max = r.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
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

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func ProcessUserData(data UserData) (UserData, error) {
	if err := ValidateEmail(data.Email); err != nil {
		return UserData{}, err
	}

	normalizedUsername := NormalizeUsername(data.Username)

	if err := ValidateAge(data.Age); err != nil {
		return UserData{}, err
	}

	return UserData{
		Email:    data.Email,
		Username: normalizedUsername,
		Age:      data.Age,
	}, nil
}