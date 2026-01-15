package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Value   float64 `json:"value"`
    Active  bool    `json:"active"`
}

func processCSVFile(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ','
    reader.Comment = '#'
    reader.FieldsPerRecord = 4

    var records []Record
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("line %d: %w", lineNumber, err)
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid ID: %w", lineNumber, err)
        }

        name := row[1]

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid value: %w", lineNumber, err)
        }

        active, err := strconv.ParseBool(row[3])
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid active flag: %w", lineNumber, err)
        }

        records = append(records, Record{
            ID:     id,
            Name:   name,
            Value:  value,
            Active: active,
        })
    }

    return records, nil
}

func generateJSONOutput(records []Record, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")

    if err := encoder.Encode(records); err != nil {
        return fmt.Errorf("failed to encode JSON: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.json>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    records, err := processCSVFile(inputFile)
    if err != nil {
        fmt.Printf("Error processing CSV: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records\n", len(records))

    if err := generateJSONOutput(records, outputFile); err != nil {
        fmt.Printf("Error generating JSON: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Output written to %s\n", outputFile)
}