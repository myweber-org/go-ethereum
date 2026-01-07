
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID      string
    Name    string
    Email   string
    Active  string
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
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

        if lineNumber == 1 {
            continue
        }

        if len(row) < 4 {
            return nil, fmt.Errorf("invalid column count at line %d", lineNumber)
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if record.ID == "" || record.Name == "" {
            return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
        }

        if !strings.Contains(record.Email, "@") {
            return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
        }

        records = append(records, record)
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) []string {
    var errors []string
    emailSet := make(map[string]bool)

    for i, record := range records {
        if record.Active != "true" && record.Active != "false" {
            errors = append(errors, fmt.Sprintf("record %d: invalid active status", i+1))
        }

        if emailSet[record.Email] {
            errors = append(errors, fmt.Sprintf("record %d: duplicate email detected", i+1))
        }
        emailSet[record.Email] = true
    }

    return errors
}

func GenerateReport(records []DataRecord) {
    activeCount := 0
    for _, record := range records {
        if record.Active == "true" {
            activeCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Active records: %d\n", activeCount)
    fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}package main

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

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) != 3 {
            return nil, errors.New("invalid csv format")
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, err
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, err
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
    if len(records) == 0 {
        return errors.New("no records to validate")
    }

    for _, record := range records {
        if record.ID <= 0 {
            return errors.New("invalid id value")
        }
        if record.Name == "" {
            return errors.New("empty name field")
        }
        if record.Value < 0 {
            return errors.New("negative value not allowed")
        }
    }

    return nil
}

func CalculateTotal(records []DataRecord) float64 {
    total := 0.0
    for _, record := range records {
        total += record.Value
    }
    return total
}