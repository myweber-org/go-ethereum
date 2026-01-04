
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
    ID    int
    Name  string
    Value float64
    Valid bool
}

type DataProcessor struct {
    records []DataRecord
    stats   struct {
        total    int
        valid    int
        sumValue float64
    }
}

func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        records: make([]DataRecord, 0),
    }
}

func (dp *DataProcessor) ProcessCSVFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    lineNumber := 0
    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if lineNumber == 1 {
            continue
        }

        record, err := dp.parseRow(row, lineNumber)
        if err != nil {
            fmt.Printf("Warning line %d: %v\n", lineNumber, err)
            continue
        }

        dp.records = append(dp.records, record)
        dp.updateStats(record)
    }

    return nil
}

func (dp *DataProcessor) parseRow(row []string, line int) (DataRecord, error) {
    if len(row) != 4 {
        return DataRecord{}, errors.New("invalid column count")
    }

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid ID format: %w", err)
    }

    name := strings.TrimSpace(row[1])
    if name == "" {
        return DataRecord{}, errors.New("name cannot be empty")
    }

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return DataRecord{}, fmt.Errorf("invalid value format: %w", err)
    }

    validStr := strings.ToLower(strings.TrimSpace(row[3]))
    valid := validStr == "true" || validStr == "1" || validStr == "yes"

    return DataRecord{
        ID:    id,
        Name:  name,
        Value: value,
        Valid: valid,
    }, nil
}

func (dp *DataProcessor) updateStats(record DataRecord) {
    dp.stats.total++
    if record.Valid {
        dp.stats.valid++
        dp.stats.sumValue += record.Value
    }
}

func (dp *DataProcessor) GetStatistics() map[string]interface{} {
    avgValue := 0.0
    if dp.stats.valid > 0 {
        avgValue = dp.stats.sumValue / float64(dp.stats.valid)
    }

    return map[string]interface{}{
        "total_records":    dp.stats.total,
        "valid_records":    dp.stats.valid,
        "invalid_records":  dp.stats.total - dp.stats.valid,
        "average_value":    fmt.Sprintf("%.2f", avgValue),
        "total_value_sum":  dp.stats.sumValue,
    }
}

func (dp *DataProcessor) FilterValidRecords() []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range dp.records {
        if record.Valid {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func (dp *DataProcessor) ExportValidCSV(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    header := []string{"ID", "Name", "Value", "Valid"}
    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    validRecords := dp.FilterValidRecords()
    for _, record := range validRecords {
        row := []string{
            strconv.Itoa(record.ID),
            record.Name,
            fmt.Sprintf("%.2f", record.Value),
            strconv.FormatBool(record.Valid),
        }
        if err := writer.Write(row); err != nil {
            return fmt.Errorf("failed to write record: %w", err)
        }
    }

    return nil
}

func main() {
    processor := NewDataProcessor()

    if err := processor.ProcessCSVFile("input.csv"); err != nil {
        fmt.Printf("Processing error: %v\n", err)
        return
    }

    stats := processor.GetStatistics()
    fmt.Println("Processing Statistics:")
    for key, value := range stats {
        fmt.Printf("%s: %v\n", key, value)
    }

    if err := processor.ExportValidCSV("valid_output.csv"); err != nil {
        fmt.Printf("Export error: %v\n", err)
        return
    }

    fmt.Println("Processing completed successfully")
}