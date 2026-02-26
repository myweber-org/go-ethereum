
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

func processCSV(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
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
            return nil, fmt.Errorf("invalid columns at line %d", lineNum)
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
        }

        value, err := strconv.ParseFloat(line[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  line[1],
            Value: value,
        })
    }

    return records, nil
}

func calculateStats(records []Record) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, r := range records {
        sum += r.Value
    }

    average := sum / float64(len(records))

    var variance float64
    for _, r := range records {
        diff := r.Value - average
        variance += diff * diff
    }
    variance /= float64(len(records))

    return average, variance
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        os.Exit(1)
    }

    records, err := processCSV(os.Args[1])
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    avg, var := calculateStats(records)
    fmt.Printf("Processed %d records\n", len(records))
    fmt.Printf("Average: %.2f\n", avg)
    fmt.Printf("Variance: %.2f\n", var)
}package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) (bool, []string) {
	var errors []string

	if strings.TrimSpace(data.Username) == "" {
		errors = append(errors, "username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		errors = append(errors, "email must contain @ symbol")
	}
	if data.Age < 0 || data.Age > 150 {
		errors = append(errors, "age must be between 0 and 150")
	}

	return len(errors) == 0, errors
}

func TransformUsername(data *UserData) {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
	user := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	TransformUsername(&user)

	if valid, errs := ValidateUserData(user); !valid {
		return user, fmt.Errorf("validation failed: %v", strings.Join(errs, "; "))
	}

	return user, nil
}

func main() {
	user, err := ProcessUserInput("  JohnDoe  ", "john@example.com", 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}