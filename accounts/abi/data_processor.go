
package main

import (
	"regexp"
	"strings"
)

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validPattern.MatchString(username)
}

func SanitizeInput(input string) string {
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type CSVProcessor struct {
    filePath string
    delimiter rune
}

func NewCSVProcessor(path string) *CSVProcessor {
    return &CSVProcessor{
        filePath:  path,
        delimiter: ',',
    }
}

func (p *CSVProcessor) SetDelimiter(d rune) {
    p.delimiter = d
}

func (p *CSVProcessor) Validate() error {
    file, err := os.Open(p.filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = p.delimiter
    reader.FieldsPerRecord = -1

    lineNumber := 0
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("line %d: %w", lineNumber, err)
        }

        for i, field := range record {
            if strings.TrimSpace(field) == "" {
                return fmt.Errorf("line %d, column %d: empty field", lineNumber, i+1)
            }
        }
        lineNumber++
    }

    if lineNumber == 0 {
        return fmt.Errorf("file is empty")
    }

    return nil
}

func (p *CSVProcessor) Clean() ([]string, error) {
    file, err := os.Open(p.filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = p.delimiter
    reader.FieldsPerRecord = -1

    var cleanedRows []string
    lineNumber := 0

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            continue
        }

        var cleanedFields []string
        for _, field := range record {
            cleaned := strings.TrimSpace(field)
            cleaned = strings.ToLower(cleaned)
            cleanedFields = append(cleanedFields, cleaned)
        }

        if len(cleanedFields) > 0 {
            row := strings.Join(cleanedFields, string(p.delimiter))
            cleanedRows = append(cleanedRows, row)
        }
        lineNumber++
    }

    return cleanedRows, nil
}

func main() {
    processor := NewCSVProcessor("data.csv")
    
    if err := processor.Validate(); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
        return
    }

    cleanedData, err := processor.Clean()
    if err != nil {
        fmt.Printf("Cleaning failed: %v\n", err)
        return
    }

    for _, row := range cleanedData {
        fmt.Println(row)
    }
}
package data

import (
	"errors"
	"strings"
	"time"
)

type Record struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(r Record) error {
	if r.ID == "" {
		return errors.New("record ID cannot be empty")
	}
	if r.Value < 0 {
		return errors.New("record value cannot be negative")
	}
	if r.Timestamp.IsZero() {
		return errors.New("record timestamp must be set")
	}
	return nil
}

func TransformRecord(r Record, multiplier float64) Record {
	return Record{
		ID:        strings.ToUpper(r.ID),
		Value:     r.Value * multiplier,
		Timestamp: r.Timestamp.UTC(),
		Tags:      append([]string{"processed"}, r.Tags...),
	}
}

func FilterRecords(records []Record, minValue float64) []Record {
	var filtered []Record
	for _, r := range records {
		if r.Value >= minValue {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func CalculateAverage(records []Record) float64 {
	if len(records) == 0 {
		return 0
	}
	
	var sum float64
	for _, r := range records {
		sum += r.Value
	}
	return sum / float64(len(records))
}