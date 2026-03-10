
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

type Record struct {
    ID      int
    Name    string
    Value   float64
    Active  bool
}

func parseCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := []Record{}
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 0 {
            lineNumber++
            continue
        }

        record, err := parseRecord(line, lineNumber)
        if err != nil {
            return nil, err
        }

        records = append(records, record)
        lineNumber++
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func parseRecord(fields []string, lineNum int) (Record, error) {
    if len(fields) != 4 {
        return Record{}, fmt.Errorf("invalid field count at line %d: expected 4, got %d", lineNum, len(fields))
    }

    id, err := strconv.Atoi(strings.TrimSpace(fields[0]))
    if err != nil {
        return Record{}, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
    }

    name := strings.TrimSpace(fields[1])
    if name == "" {
        return Record{}, fmt.Errorf("empty name at line %d", lineNum)
    }

    value, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
    if err != nil {
        return Record{}, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
    }

    active, err := strconv.ParseBool(strings.TrimSpace(fields[3]))
    if err != nil {
        return Record{}, fmt.Errorf("invalid active flag at line %d: %w", lineNum, err)
    }

    return Record{
        ID:     id,
        Name:   name,
        Value:  value,
        Active: active,
    }, nil
}

func calculateStats(records []Record) (float64, float64, int) {
    if len(records) == 0 {
        return 0, 0, 0
    }

    var sum float64
    var activeCount int
    var minValue float64 = records[0].Value
    var maxValue float64 = records[0].Value

    for _, record := range records {
        sum += record.Value
        if record.Active {
            activeCount++
        }
        if record.Value < minValue {
            minValue = record.Value
        }
        if record.Value > maxValue {
            maxValue = record.Value
        }
    }

    average := sum / float64(len(records))
    return average, maxValue - minValue, activeCount
}

func filterRecords(records []Record, predicate func(Record) bool) []Record {
    filtered := []Record{}
    for _, record := range records {
        if predicate(record) {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        os.Exit(1)
    }

    filename := os.Args[1]
    records, err := parseCSVFile(filename)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    avgValue, valueRange, activeCount := calculateStats(records)
    fmt.Printf("Processed %d records\n", len(records))
    fmt.Printf("Average value: %.2f\n", avgValue)
    fmt.Printf("Value range: %.2f\n", valueRange)
    fmt.Printf("Active records: %d\n", activeCount)

    activeRecords := filterRecords(records, func(r Record) bool {
        return r.Active && r.Value > 50.0
    })
    fmt.Printf("High-value active records: %d\n", len(activeRecords))
}