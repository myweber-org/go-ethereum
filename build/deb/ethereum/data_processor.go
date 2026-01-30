
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

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processCSVFile(inputPath string, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()

    csvReader := csv.NewReader(inputFile)
    csvWriter := csv.NewWriter(outputFile)
    defer csvWriter.Flush()

    headers, err := csvReader.Read()
    if err != nil {
        return fmt.Errorf("failed to read headers: %w", err)
    }

    if err := csvWriter.Write(headers); err != nil {
        return fmt.Errorf("failed to write headers: %w", err)
    }

    recordCount := 0
    validCount := 0

    for {
        row, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read row: %w", err)
        }

        recordCount++
        if len(row) < 4 {
            continue
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if record.ID == "" || record.Name == "" {
            continue
        }

        if !validateEmail(record.Email) {
            continue
        }

        if record.Active != "true" && record.Active != "false" {
            continue
        }

        outputRow := []string{
            record.ID,
            strings.ToUpper(record.Name),
            strings.ToLower(record.Email),
            record.Active,
        }

        if err := csvWriter.Write(outputRow); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }

        validCount++
    }

    fmt.Printf("Processed %d records, %d valid records written to %s\n", 
        recordCount, validCount, outputPath)
    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run data_processor.go <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }
}