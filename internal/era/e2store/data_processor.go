
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type DataRecord struct {
	ID    int
	Name  string
	Value float64
}

func parseCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
	lineNum := 0

	for {
		lineNum++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error at line %d: %w", lineNum, err)
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNum, len(row))
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", lineNum, err)
		}

		name := row[1]
		if name == "" {
			return nil, fmt.Errorf("empty name at line %d", lineNum)
		}

		value, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value at line %d: %w", lineNum, err)
		}

		records = append(records, DataRecord{
			ID:    id,
			Name:  name,
			Value: value,
		})
	}

	return records, nil
}

func calculateStats(records []DataRecord) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var max float64
	count := len(records)

	for i, record := range records {
		sum += record.Value
		if i == 0 || record.Value > max {
			max = record.Value
		}
	}

	average := sum / float64(count)
	return average, max, count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	records, err := parseCSVFile(filename)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	avg, max, count := calculateStats(records)
	fmt.Printf("Processed %d records\n", count)
	fmt.Printf("Average value: %.2f\n", avg)
	fmt.Printf("Maximum value: %.2f\n", max)
}
package main

import (
    "errors"
    "fmt"
    "strings"
    "time"
)

type DataRecord struct {
    ID        string
    Timestamp time.Time
    Value     float64
    Tags      []string
}

func ValidateRecord(record DataRecord) error {
    if record.ID == "" {
        return errors.New("record ID cannot be empty")
    }
    if record.Value < 0 {
        return errors.New("record value cannot be negative")
    }
    if record.Timestamp.After(time.Now()) {
        return errors.New("record timestamp cannot be in the future")
    }
    return nil
}

func TransformTags(tags []string) []string {
    var transformed []string
    for _, tag := range tags {
        trimmed := strings.TrimSpace(tag)
        if trimmed != "" {
            transformed = append(transformed, strings.ToLower(trimmed))
        }
    }
    return transformed
}

func CalculateAverage(records []DataRecord) (float64, error) {
    if len(records) == 0 {
        return 0, errors.New("cannot calculate average for empty record set")
    }

    var sum float64
    validCount := 0

    for _, record := range records {
        if err := ValidateRecord(record); err == nil {
            sum += record.Value
            validCount++
        }
    }

    if validCount == 0 {
        return 0, errors.New("no valid records found for average calculation")
    }

    return sum / float64(validCount), nil
}

func FilterByTag(records []DataRecord, targetTag string) []DataRecord {
    var filtered []DataRecord
    for _, record := range records {
        for _, tag := range record.Tags {
            if strings.EqualFold(tag, targetTag) {
                filtered = append(filtered, record)
                break
            }
        }
    }
    return filtered
}

func ProcessDataBatch(records []DataRecord) ([]DataRecord, error) {
    var processed []DataRecord
    var validationErrors []string

    for i, record := range records {
        if err := ValidateRecord(record); err != nil {
            validationErrors = append(validationErrors, fmt.Sprintf("record %d: %v", i, err))
            continue
        }

        processedRecord := DataRecord{
            ID:        record.ID,
            Timestamp: record.Timestamp,
            Value:     record.Value,
            Tags:      TransformTags(record.Tags),
        }
        processed = append(processed, processedRecord)
    }

    if len(validationErrors) > 0 && len(processed) == 0 {
        return nil, fmt.Errorf("all records failed validation: %v", strings.Join(validationErrors, "; "))
    }

    return processed, nil
}