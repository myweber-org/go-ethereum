
package data_processor

import (
	"regexp"
	"strings"
	"unicode"
)

type Processor struct {
	stripSpaces   bool
	removeSpecial bool
	toLowercase   bool
}

func NewProcessor(stripSpaces, removeSpecial, toLowercase bool) *Processor {
	return &Processor{
		stripSpaces:   stripSpaces,
		removeSpecial: removeSpecial,
		toLowercase:   toLowercase,
	}
}

func (p *Processor) CleanString(input string) string {
	result := input

	if p.stripSpaces {
		result = strings.TrimSpace(result)
		result = strings.Join(strings.Fields(result), " ")
	}

	if p.removeSpecial {
		reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
		result = reg.ReplaceAllString(result, "")
	}

	if p.toLowercase {
		result = strings.ToLower(result)
	}

	return result
}

func (p *Processor) NormalizeWhitespace(input string) string {
	var builder strings.Builder
	prevSpace := false

	for _, r := range input {
		if unicode.IsSpace(r) {
			if !prevSpace {
				builder.WriteRune(' ')
				prevSpace = true
			}
		} else {
			builder.WriteRune(r)
			prevSpace = false
		}
	}

	return strings.TrimSpace(builder.String())
}

func (p *Processor) ProcessBatch(inputs []string) []string {
	results := make([]string, len(inputs))
	for i, input := range inputs {
		results[i] = p.CleanString(input)
	}
	return results
}
package main

import (
    "encoding/csv"
    "errors"
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

func ProcessCSVFile(filename string) ([]DataRecord, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records := make([]DataRecord, 0)

    lineNumber := 0
    for {
        lineNumber++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("csv read error at line %d: %w", lineNumber, err)
        }

        if len(row) != 3 {
            return nil, fmt.Errorf("invalid column count at line %d: expected 3, got %d", lineNumber, len(row))
        }

        id, err := strconv.Atoi(row[0])
        if err != nil {
            return nil, fmt.Errorf("invalid ID at line %d: %w", lineNumber, err)
        }

        name := row[1]
        if name == "" {
            return nil, fmt.Errorf("empty name at line %d", lineNumber)
        }

        value, err := strconv.ParseFloat(row[2], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid value at line %d: %w", lineNumber, err)
        }

        records = append(records, DataRecord{
            ID:    id,
            Name:  name,
            Value: value,
        })
    }

    if len(records) == 0 {
        return nil, errors.New("no valid records found in file")
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) error {
    idSet := make(map[int]bool)
    for _, record := range records {
        if record.ID <= 0 {
            return fmt.Errorf("invalid ID %d: must be positive", record.ID)
        }
        if idSet[record.ID] {
            return fmt.Errorf("duplicate ID %d found", record.ID)
        }
        idSet[record.ID] = true

        if record.Value < 0 {
            return fmt.Errorf("negative value %f for record ID %d", record.Value, record.ID)
        }
    }
    return nil
}

func CalculateStatistics(records []DataRecord) (float64, float64, int) {
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