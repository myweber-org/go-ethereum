
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
    ID    int
    Name  string
    Email string
    Age   int
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

    headers, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read headers: %w", err)
    }

    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("failed to write headers: %w", err)
    }

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            continue
        }

        if len(row) != 4 {
            continue
        }

        id, err := strconv.Atoi(strings.TrimSpace(row[0]))
        if err != nil || id <= 0 {
            continue
        }

        name := strings.TrimSpace(row[1])
        if name == "" {
            continue
        }

        email := strings.TrimSpace(row[2])
        if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
            continue
        }

        age, err := strconv.Atoi(strings.TrimSpace(row[3]))
        if err != nil || age < 0 || age > 120 {
            continue
        }

        cleanedRow := []string{
            strconv.Itoa(id),
            name,
            email,
            strconv.Itoa(age),
        }

        if err := writer.Write(cleanedRow); err != nil {
            return fmt.Errorf("failed to write row: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputPath := os.Args[1]
    outputPath := os.Args[2]

    if err := cleanCSV(inputPath, outputPath); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}