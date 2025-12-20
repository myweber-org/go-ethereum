
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID    string
    Name  string
    Email string
    Valid bool
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

        if len(row) < 3 {
            continue
        }

        record := DataRecord{
            ID:    strings.TrimSpace(row[0]),
            Name:  strings.TrimSpace(row[1]),
            Email: strings.TrimSpace(row[2]),
            Valid: validateRecord(strings.TrimSpace(row[0]), strings.TrimSpace(row[2])),
        }

        records = append(records, record)
    }

    return records, nil
}

func validateRecord(id, email string) bool {
    if id == "" || email == "" {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FilterValidRecords(records []DataRecord) []DataRecord {
    var validRecords []DataRecord
    for _, record := range records {
        if record.Valid {
            validRecords = append(validRecords, record)
        }
    }
    return validRecords
}

func GenerateReport(records []DataRecord) {
    validCount := 0
    for _, record := range records {
        if record.Valid {
            validCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Valid records: %d\n", validCount)
    fmt.Printf("Invalid records: %d\n", len(records)-validCount)
    
    if validCount > 0 {
        fmt.Println("\nValid records:")
        for _, record := range records {
            if record.Valid {
                fmt.Printf("  ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
            }
        }
    }
}