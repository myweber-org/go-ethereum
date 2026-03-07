package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

func cleanCSV(inputPath, outputPath string) error {
    inFile, err := os.Open(inputPath)
    if err != nil {
        return fmt.Errorf("failed to open input file: %w", err)
    }
    defer inFile.Close()

    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("error reading CSV: %w", err)
        }

        cleaned := make([]string, len(record))
        for i, field := range record {
            cleaned[i] = strings.TrimSpace(field)
        }

        if err := writer.Write(cleaned); err != nil {
            return fmt.Errorf("error writing CSV: %w", err)
        }
    }

    return nil
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    err := cleanCSV(os.Args[1], os.Args[2])
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("CSV cleaning completed successfully")
}
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
    Score float64
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

    headers := []string{"ID", "Name", "Email", "Score"}
    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("failed to write headers: %w", err)
    }

    lineNum := 0
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("warning: line %d: %v\n", lineNum, err)
            continue
        }

        if len(row) != 4 {
            fmt.Printf("warning: line %d: invalid column count\n", lineNum)
            continue
        }

        record, err := parseRecord(row)
        if err != nil {
            fmt.Printf("warning: line %d: %v\n", lineNum, err)
            continue
        }

        if !validateRecord(record) {
            fmt.Printf("warning: line %d: validation failed\n", lineNum)
            continue
        }

        outputRow := []string{
            strconv.Itoa(record.ID),
            strings.TrimSpace(record.Name),
            strings.ToLower(strings.TrimSpace(record.Email)),
            fmt.Sprintf("%.2f", record.Score),
        }

        if err := writer.Write(outputRow); err != nil {
            return fmt.Errorf("failed to write record: %w", err)
        }
    }

    return nil
}

func parseRecord(row []string) (Record, error) {
    var record Record

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, fmt.Errorf("invalid ID: %w", err)
    }
    record.ID = id

    record.Name = row[1]

    record.Email = row[2]

    score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
    if err != nil {
        return record, fmt.Errorf("invalid score: %w", err)
    }
    record.Score = score

    return record, nil
}

func validateRecord(r Record) bool {
    if r.ID <= 0 {
        return false
    }
    if strings.TrimSpace(r.Name) == "" {
        return false
    }
    if !strings.Contains(r.Email, "@") {
        return false
    }
    if r.Score < 0 || r.Score > 100 {
        return false
    }
    return true
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputPath := os.Args[1]
    outputPath := os.Args[2]

    if err := cleanCSV(inputPath, outputPath); err != nil {
        fmt.Printf("error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("data cleaning completed successfully")
}