
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

func ValidateEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

func SanitizeUsername(username string) string {
    return strings.TrimSpace(username)
}

func TransformUserData(rawData []byte) (*UserData, error) {
    var data UserData
    err := json.Unmarshal(rawData, &data)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    if !ValidateEmail(data.Email) {
        return nil, fmt.Errorf("invalid email format: %s", data.Email)
    }

    data.Username = SanitizeUsername(data.Username)

    if data.Age < 0 || data.Age > 150 {
        return nil, fmt.Errorf("age out of valid range: %d", data.Age)
    }

    return &data, nil
}

func main() {
    rawJSON := []byte(`{"email":"test@example.com","username":"  user123  ","age":25}`)
    processedData, err := TransformUserData(rawJSON)
    if err != nil {
        fmt.Printf("Error processing data: %v\n", err)
        return
    }
    fmt.Printf("Processed Data: %+v\n", processedData)
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
			return nil, fmt.Errorf("invalid column count at line %d", lineNum)
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		name := line[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNum)
		}

		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
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
}