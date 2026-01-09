
package main

import (
    "encoding/csv"
    "errors"
    "io"
    "os"
    "strconv"
    "strings"
)

type DataRecord struct {
    ID    int
    Name  string
    Value float64
    Valid bool
}

func ParseCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    for line := 1; ; line++ {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        if len(row) < 4 {
            continue
        }

        record, err := parseRow(row, line)
        if err != nil {
            continue
        }

        records = append(records, record)
    }

    return records, nil
}

func parseRow(row []string, line int) (DataRecord, error) {
    var record DataRecord

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil {
        return record, errors.New("invalid id")
    }
    record.ID = id

    name := strings.TrimSpace(row[1])
    if name == "" {
        return record, errors.New("empty name")
    }
    record.Name = name

    value, err := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)
    if err != nil {
        return record, errors.New("invalid value")
    }
    record.Value = value

    valid := strings.ToLower(strings.TrimSpace(row[3]))
    record.Valid = valid == "true" || valid == "yes" || valid == "1"

    return record, nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
    filtered := make([]DataRecord, 0)
    for _, record := range records {
        if record.Valid {
            filtered = append(filtered, record)
        }
    }
    return filtered
}

func CalculateAverage(records []DataRecord) float64 {
    if len(records) == 0 {
        return 0
    }

    total := 0.0
    count := 0
    for _, record := range records {
        if record.Valid {
            total += record.Value
            count++
        }
    }

    if count == 0 {
        return 0
    }
    return total / float64(count)
}