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
    Score   float64
    Valid   bool
}

func cleanCSV(inputPath, outputPath string) error {
    inFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inFile.Close()

    outFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outFile.Close()

    reader := csv.NewReader(inFile)
    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    header, err := reader.Read()
    if err != nil {
        return err
    }
    header = append(header, "Valid")
    writer.Write(header)

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            continue
        }

        record := parseRecord(row)
        record.Valid = validateRecord(record)

        outputRow := []string{
            strconv.Itoa(record.ID),
            strings.TrimSpace(record.Name),
            strings.ToLower(strings.TrimSpace(record.Email)),
            fmt.Sprintf("%.2f", record.Score),
            strconv.FormatBool(record.Valid),
        }
        writer.Write(outputRow)
    }
    return nil
}

func parseRecord(row []string) Record {
    id, _ := strconv.Atoi(row[0])
    score, _ := strconv.ParseFloat(row[3], 64)
    return Record{
        ID:    id,
        Name:  row[1],
        Email: row[2],
        Score: score,
    }
}

func validateRecord(r Record) bool {
    if r.ID <= 0 {
        return false
    }
    if len(r.Name) == 0 || len(r.Name) > 100 {
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
    err := cleanCSV("input.csv", "cleaned.csv")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Data cleaning completed successfully")
}