package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
)

func processCSVFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    lineCount := 0

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("csv read error: %w", err)
        }

        lineCount++
        if lineCount == 1 {
            fmt.Println("Headers:", record)
            continue
        }

        fmt.Printf("Row %d: %v\n", lineCount-1, record)
    }

    fmt.Printf("Total records processed: %d\n", lineCount-1)
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run file_processor.go <csv_file>")
        os.Exit(1)
    }

    if err := processCSVFile(os.Args[1]); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}