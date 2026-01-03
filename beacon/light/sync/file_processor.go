package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
)

type Record struct {
    ID    int
    Name  string
    Value float64
}

type FileProcessor struct {
    logger *log.Logger
}

func NewFileProcessor() *FileProcessor {
    return &FileProcessor{
        logger: log.New(os.Stderr, "PROCESSOR: ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}

func (fp *FileProcessor) ValidateRecord(rec Record) error {
    if rec.ID <= 0 {
        return errors.New("invalid ID")
    }
    if rec.Name == "" {
        return errors.New("empty name")
    }
    if rec.Value < 0 {
        return errors.New("negative value")
    }
    return nil
}

func (fp *FileProcessor) ProcessCSV(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNumber := 0

    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fp.logger.Printf("line %d: CSV read error: %v", lineNumber, err)
            continue
        }

        if len(row) != 3 {
            fp.logger.Printf("line %d: invalid column count: %d", lineNumber, len(row))
            continue
        }

        id, err1 := strconv.Atoi(row[0])
        value, err2 := strconv.ParseFloat(row[2], 64)
        if err1 != nil || err2 != nil {
            fp.logger.Printf("line %d: parse error - ID:%v, Value:%v", lineNumber, err1, err2)
            continue
        }

        record := Record{
            ID:    id,
            Name:  row[1],
            Value: value,
        }

        if err := fp.ValidateRecord(record); err != nil {
            fp.logger.Printf("line %d: validation failed: %v", lineNumber, err)
            continue
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records processed")
    }

    return records, nil
}

func (fp *FileProcessor) GenerateSummary(records []Record) {
    total := len(records)
    var sum float64
    for _, r := range records {
        sum += r.Value
    }
    average := sum / float64(total)

    fp.logger.Printf("Processed %d valid records", total)
    fp.logger.Printf("Total value: %.2f", sum)
    fp.logger.Printf("Average value: %.2f", average)
}

func main() {
    processor := NewFileProcessor()
    
    if len(os.Args) < 2 {
        processor.logger.Fatal("Usage: file_processor <csv_filename>")
    }

    filename := os.Args[1]
    records, err := processor.ProcessCSV(filename)
    if err != nil {
        processor.logger.Fatal(err)
    }

    processor.GenerateSummary(records)
    fmt.Printf("Successfully processed %d records from %s\n", len(records), filename)
}