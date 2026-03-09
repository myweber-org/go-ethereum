
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
)

type UserData struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func ValidateJSON(input string) (*UserData, error) {
    if strings.TrimSpace(input) == "" {
        return nil, errors.New("empty input string")
    }

    var data UserData
    err := json.Unmarshal([]byte(input), &data)
    if err != nil {
        return nil, fmt.Errorf("JSON parsing failed: %w", err)
    }

    if data.Name == "" {
        return nil, errors.New("name field is required")
    }
    if data.Email == "" {
        return nil, errors.New("email field is required")
    }
    if data.Age <= 0 {
        return nil, errors.New("age must be a positive integer")
    }

    return &data, nil
}

func ProcessUserData(jsonStr string) {
    userData, err := ValidateJSON(jsonStr)
    if err != nil {
        fmt.Printf("Validation error: %v\n", err)
        return
    }

    fmt.Printf("Valid user data: Name=%s, Email=%s, Age=%d\n",
        userData.Name, userData.Email, userData.Age)
}

func main() {
    testData := `{"name":"John Doe","email":"john@example.com","age":30}`
    ProcessUserData(testData)

    invalidData := `{"name":"","email":"test@example.com","age":25}`
    ProcessUserData(invalidData)
}
package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for lineNumber := 1; ; lineNumber++ {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        if value < 0 {
            return nil, fmt.Errorf("negative value at line %d: %f", lineNumber, value)
        }

        records = append(records, DataRecord{
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

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
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

    avg, max, count := CalculateStatistics(records)
    fmt.Printf("Processed %d records\n", count)
    fmt.Printf("Average value: %.2f\n", avg)
    fmt.Printf("Maximum value: %.2f\n", max)
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
	reader.TrimLeadingSpace = true

	var records []Record
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

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNumber)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func CalculateStats(records []Record) (float64, float64) {
	if len(records) == 0 {
		return 0, 0
	}

	var sum float64
	var max float64 = records[0].Value

	for _, record := range records {
		sum += record.Value
		if record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(len(records))
	return average, max
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := ProcessCSV(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	average, max := CalculateStats(records)
	fmt.Printf("Processed %d records\n", len(records))
	fmt.Printf("Average value: %.2f\n", average)
	fmt.Printf("Maximum value: %.2f\n", max)
}