package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    ID      int
    Name    string
    Email   string
    Active  bool
    Score   float64
}

func cleanCSV(inputPath, outputPath string) error {
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

    reader := csv.NewReader(inputFile)
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read header: %w", err)
    }

    if err := writer.Write(header); err != nil {
        return fmt.Errorf("failed to write header: %w", err)
    }

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("line %d: skipped due to read error: %v\n", lineNum, err)
            continue
        }

        if len(row) != 5 {
            fmt.Printf("line %d: skipped due to column count mismatch\n", lineNum)
            continue
        }

        record, err := parseRecord(row)
        if err != nil {
            fmt.Printf("line %d: skipped due to validation error: %v\n", lineNum, err)
            continue
        }

        cleaned := cleanRecord(record)
        outputRow := []string{
            strconv.Itoa(cleaned.ID),
            cleaned.Name,
            cleaned.Email,
            strconv.FormatBool(cleaned.Active),
            fmt.Sprintf("%.2f", cleaned.Score),
        }

        if err := writer.Write(outputRow); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }
    }

    return nil
}

func parseRecord(row []string) (Record, error) {
    var r Record
    var err error

    if r.ID, err = strconv.Atoi(strings.TrimSpace(row[0])); err != nil {
        return r, fmt.Errorf("invalid ID: %w", err)
    }

    r.Name = strings.TrimSpace(row[1])
    if r.Name == "" {
        return r, fmt.Errorf("name cannot be empty")
    }

    r.Email = strings.TrimSpace(row[2])
    if !strings.Contains(r.Email, "@") {
        return r, fmt.Errorf("invalid email format")
    }

    activeStr := strings.ToLower(strings.TrimSpace(row[3]))
    r.Active = activeStr == "true" || activeStr == "1" || activeStr == "yes"

    if r.Score, err = strconv.ParseFloat(strings.TrimSpace(row[4]), 64); err != nil {
        return r, fmt.Errorf("invalid score: %w", err)
    }

    return r, nil
}

func cleanRecord(r Record) Record {
    r.Name = strings.Title(strings.ToLower(r.Name))
    r.Email = strings.ToLower(r.Email)
    if r.Score < 0 {
        r.Score = 0
    } else if r.Score > 100 {
        r.Score = 100
    }
    return r
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := cleanCSV(inputFile, outputFile); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}