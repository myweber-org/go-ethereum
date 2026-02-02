
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
    ID      int
    Name    string
    Value   float64
    Active  bool
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

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, fmt.Errorf("empty name at line %d", lineNumber)
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }
    record.Value = value

    active, err := strconv.ParseBool(strings.TrimSpace(row[3]))
    if err != nil {
        return record, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
    }
    record.Active = active

    return record, nil
}

func ValidateRecords(records []DataRecord) error {
    idSet := make(map[int]bool)

    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid record ID: %d (must be positive)", record.ID)
        }

        if idSet[record.ID] {
            return fmt.Errorf("duplicate ID found: %d", record.ID)
        }
        idSet[record.ID] = true

        if record.Value < 0 {
            return fmt.Errorf("negative value for record ID %d: %f", record.ID, record.Value)
        }
    }

    return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var activeCount int
    minValue := records[0].Value

    for _, record := range records {
        sum += record.Value
        if record.Value < minValue {
            minValue = record.Value
        }
        if record.Active {
            activeCount++
        }
    }

    average := sum / float64(len(records))
    return average, minValue, activeCount
}package main

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
    var records []Record
    line := 0

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        line++
        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d", line)
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", line, err)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", line, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        })
    }

    return records, nil
}

func ValidateRecords(records []Record) error {
    seen := make(map[int]bool)
    for _, r := range records {
        if r.ID <= 0 {
            return fmt.Errorf("invalid ID %d", r.ID)
        }
        if r.Name == "" {
            return fmt.Errorf("empty name for ID %d", r.ID)
        }
        if r.Value < 0 {
            return fmt.Errorf("negative value for ID %d", r.ID)
        }
        if seen[r.ID] {
            return fmt.Errorf("duplicate ID %d", r.ID)
        }
        seen[r.ID] = true
    }
    return nil
}

func CalculateStats(records []Record) (float64, float64) {
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
    stdDev := variance / float64(len(records))

    return average, stdDev
}