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
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Value float64 `json:"value"`
}

func processCSVFile(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ','
    reader.TrimLeadingSpace = true

    var records []Record
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("line %d: %v", lineNumber, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("line %d: expected 3 columns, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid ID: %v", lineNumber, err)
        }

        name := row[1]

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("line %d: invalid value: %v", lineNumber, err)
        }

        records = append(records, Record{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    return records, nil
}

func writeJSONFile(records []Record, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(records)
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

    err = writeJSONFile(records, outputFile)
    if err != nil {
        fmt.Printf("Error writing JSON: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records to %s\n", len(records), outputFile)
}