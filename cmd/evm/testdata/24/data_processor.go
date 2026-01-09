
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

type DataProcessor struct {
    records []DataRecord
}

func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        records: make([]DataRecord, 0),
    }
}

func (dp *DataProcessor) LoadCSV(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    lineNumber := 0

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 0 {
            lineNumber++
            continue
        }

        dataRecord, err := parseRecord(record, lineNumber)
        if err != nil {
            return err
        }

        dp.records = append(dp.records, dataRecord)
        lineNumber++
    }

    if len(dp.records) == 0 {
        return errors.New("no valid records found in CSV")
    }

    return nil
}

func parseRecord(fields []string, lineNumber int) (DataRecord, error) {
    if len(fields) != 4 {
        return DataRecord{}, fmt.Errorf("invalid field count at line %d: expected 4, got %d", lineNumber, len(fields))
    }

    id, err := strconv.Atoi(fields[0])
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
    }

    name := strings.TrimSpace(fields[1])
    if name == "" {
        return DataRecord{}, fmt.Errorf("empty name at line %d", lineNumber)
    }

    value, err := strconv.ParseFloat(fields[2], 64)
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
    }

    active, err := strconv.ParseBool(fields[3])
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid active flag at line %d: %w", lineNumber, err)
    }

    return DataRecord{
        ID:     id,
        Name:   name,
        Value:  value,
        Active: active,
    }, nil
}

func (dp *DataProcessor) FilterActive() []DataRecord {
    var activeRecords []DataRecord
    for _, record := range dp.records {
        if record.Active {
            activeRecords = append(activeRecords, record)
        }
    }
    return activeRecords
}

func (dp *DataProcessor) CalculateTotal() float64 {
    var total float64
    for _, record := range dp.records {
        total += record.Value
    }
    return total
}

func (dp *DataProcessor) FindByName(name string) *DataRecord {
    for _, record := range dp.records {
        if strings.EqualFold(record.Name, name) {
            return &record
        }
    }
    return nil
}

func (dp *DataProcessor) RecordCount() int {
    return len(dp.records)
}

func main() {
    processor := NewDataProcessor()
    
    err := processor.LoadCSV("data.csv")
    if err != nil {
        fmt.Printf("Error loading CSV: %v\n", err)
        return
    }
    
    fmt.Printf("Loaded %d records\n", processor.RecordCount())
    fmt.Printf("Total value: %.2f\n", processor.CalculateTotal())
    
    activeRecords := processor.FilterActive()
    fmt.Printf("Active records: %d\n", len(activeRecords))
    
    if record := processor.FindByName("Test"); record != nil {
        fmt.Printf("Found record: ID=%d, Value=%.2f\n", record.ID, record.Value)
    }
}