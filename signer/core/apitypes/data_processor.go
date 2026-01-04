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
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Value string `json:"value"`
}

func ProcessCSVFile(inputPath string) ([]Record, error) {
    file, err := os.Open(inputPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error: %w", err)
        }

        if len(row) < 3 {
            continue
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            continue
        }

        record := Record{
            ID:    id,
            Name:  row[1],
            Value: row[2],
        }
        records = append(records, record)
    }

    return records, nil
}

func ConvertToJSON(records []Record) (string, error) {
    jsonData, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return "", fmt.Errorf("json marshal error: %w", err)
    }
    return string(jsonData), nil
}

func WriteJSONToFile(jsonData string, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    _, err = file.WriteString(jsonData)
    if err != nil {
        return fmt.Errorf("write error: %w", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: data_processor <input.csv> <output.json>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    records, err := ProcessCSVFile(inputFile)
    if err != nil {
        fmt.Printf("Error processing CSV: %v\n", err)
        os.Exit(1)
    }

    jsonData, err := ConvertToJSON(records)
    if err != nil {
        fmt.Printf("Error converting to JSON: %v\n", err)
        os.Exit(1)
    }

    err = WriteJSONToFile(jsonData, outputFile)
    if err != nil {
        fmt.Printf("Error writing JSON file: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records to %s\n", len(records), outputFile)
}