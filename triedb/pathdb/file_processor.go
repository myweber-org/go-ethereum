package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
)

type Record struct {
    ID    int
    Name  string
    Value float64
}

type FileProcessor struct {
    logger io.Writer
}

func NewFileProcessor(logger io.Writer) *FileProcessor {
    return &FileProcessor{logger: logger}
}

func (fp *FileProcessor) ProcessCSV(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]Record, 0)
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fp.logError(lineNumber, "read error", err)
            continue
        }

        if len(row) != 3 {
            fp.logError(lineNumber, "invalid column count", errors.New("expected 3 columns"))
            continue
        }

        record, err := fp.parseRow(lineNumber, row)
        if err != nil {
            fp.logError(lineNumber, "parse error", err)
            continue
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records processed")
    }

    return records, nil
}

func (fp *FileProcessor) parseRow(line int, row []string) (Record, error) {
    id, err := strconv.Atoi(row[0])
    if err != nil {
        return Record{}, fmt.Errorf("invalid ID: %w", err)
    }

    name := row[1]
    if name == "" {
        return Record{}, errors.New("name cannot be empty")
    }

    value, err := strconv.ParseFloat(row[2], 64)
    if err != nil {
        return Record{}, fmt.Errorf("invalid value: %w", err)
    }

    return Record{
        ID:    id,
        Name:  name,
        Value: value,
    }, nil
}

func (fp *FileProcessor) logError(line int, context string, err error) {
    if fp.logger != nil {
        msg := fmt.Sprintf("Line %d: %s - %v\n", line, context, err)
        fp.logger.Write([]byte(msg))
    }
}

func main() {
    processor := NewFileProcessor(os.Stderr)
    
    records, err := processor.ProcessCSV("data.csv")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Processing failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully processed %d records\n", len(records))
    for _, r := range records {
        fmt.Printf("ID: %d, Name: %s, Value: %.2f\n", r.ID, r.Name, r.Value)
    }
}