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

func cleanString(s string) string {
    return strings.TrimSpace(strings.ToLower(s))
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

    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("failed to read record: %w", err)
        }

        if len(record) < 4 {
            continue
        }

        data := DataRecord{
            ID:     cleanString(record[0]),
            Name:   cleanString(record[1]),
            Email:  cleanString(record[2]),
            Active: cleanString(record[3]),
        }

        if data.ID == "" || data.Name == "" {
            continue
        }

        if !validateEmail(data.Email) {
            continue
        }

        if data.Active != "true" && data.Active != "false" {
            continue
        }

        cleanedRecord := []string{
            data.ID,
            strings.Title(data.Name),
            data.Email,
            data.Active,
        }

        if err := csvWriter.Write(cleanedRecord); err != nil {
            return fmt.Errorf("failed to write record: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := processCSVFile(inputFile, outputFile); err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully cleaned data. Output saved to %s\n", outputFile)
}