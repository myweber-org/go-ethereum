
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
    ID        int     `json:"id"`
    Name      string  `json:"name"`
    Value     float64 `json:"value"`
    Timestamp string  `json:"timestamp"`
}

func processCSVFile(filename string) ([]Record, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    var records []Record
    lineNumber := 0

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        lineNumber++
        if lineNumber == 1 {
            continue
        }

        if len(line) != 4 {
            continue
        }

        id, err := strconv.Atoi(line[0])
        if err != nil {
            continue
        }

        value, err := strconv.ParseFloat(line[2], 64)
        if err != nil {
            continue
        }

        record := Record{
            ID:        id,
            Name:      line[1],
            Value:     value,
            Timestamp: line[3],
        }
        records = append(records, record)
    }

    return records, nil
}

func convertToJSON(records []Record) (string, error) {
    jsonData, err := json.MarshalIndent(records, "", "  ")
    if err != nil {
        return "", err
    }
    return string(jsonData), nil
}

func calculateStatistics(records []Record) (float64, float64) {
    if len(records) == 0 {
        return 0, 0
    }

    var sum float64
    for _, record := range records {
        sum += record.Value
    }
    average := sum / float64(len(records))

    var varianceSum float64
    for _, record := range records {
        diff := record.Value - average
        varianceSum += diff * diff
    }
    variance := varianceSum / float64(len(records))

    return average, variance
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file>")
        return
    }

    filename := os.Args[1]
    records, err := processCSVFile(filename)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        return
    }

    fmt.Printf("Processed %d records\n", len(records))

    avg, variance := calculateStatistics(records)
    fmt.Printf("Average value: %.2f\n", avg)
    fmt.Printf("Variance: %.2f\n", variance)

    jsonOutput, err := convertToJSON(records)
    if err != nil {
        fmt.Printf("Error converting to JSON: %v\n", err)
        return
    }

    outputFile := "output.json"
    err = os.WriteFile(outputFile, []byte(jsonOutput), 0644)
    if err != nil {
        fmt.Printf("Error writing JSON file: %v\n", err)
        return
    }

    fmt.Printf("JSON output written to %s\n", outputFile)
}