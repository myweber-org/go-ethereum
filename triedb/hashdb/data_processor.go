
package main

import (
    "encoding/csv"
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

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        records = append(records, DataRecord{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) error {
    seenIDs := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid record ID: %d", record.ID)
        }
        if record.Name == "" {
            return fmt.Errorf("empty name for record ID: %d", record.ID)
        }
        if record.Value < 0 {
            return fmt.Errorf("negative value for record ID: %d", record.ID)
        }
        if seenIDs[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        seenIDs[record.ID] = true
    }
    return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64) {
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

	data.Name = strings.TrimSpace(data.Name)
	if data.Name == "" {
		return errors.New("name cannot be empty")
	}

	data.Email = strings.TrimSpace(data.Email)
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}

	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	return nil
}

func TransformUserName(name string) string {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return "Anonymous"
	}
	return strings.Title(strings.ToLower(name))
}

func ProcessUserInput(rawName string, rawEmail string, rawAge int) (UserData, error) {
	transformedName := TransformUserName(rawName)

	user := UserData{
		Name:  transformedName,
		Email: strings.ToLower(strings.TrimSpace(rawEmail)),
		Age:   rawAge,
	}

	if err := ValidateUserData(user); err != nil {
		return UserData{}, err
	}

	return user, nil
}
package main

import (
	"fmt"
)

// FilterAndDouble filters even numbers from a slice and doubles their values.
func FilterAndDouble(numbers []int) []int {
	var result []int
	for _, num := range numbers {
		if num%2 == 0 {
			result = append(result, num*2)
		}
	}
	return result
}

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := FilterAndDouble(input)
	fmt.Println("Original:", input)
	fmt.Println("Filtered and Doubled:", output)
}