package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TrimAndLower(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "strings"
)

type DataRecord struct {
    ID      string
    Name    string
    Email   string
    Active  string
}

func ProcessCSVFile(filePath string) ([]DataRecord, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.TrimLeadingSpace = true

    var records []DataRecord
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

        if lineNumber == 1 {
            continue
        }

        if len(row) < 4 {
            return nil, fmt.Errorf("insufficient columns at line %d", lineNumber)
        }

        record := DataRecord{
            ID:     strings.TrimSpace(row[0]),
            Name:   strings.TrimSpace(row[1]),
            Email:  strings.TrimSpace(row[2]),
            Active: strings.TrimSpace(row[3]),
        }

        if record.ID == "" || record.Name == "" {
            return nil, fmt.Errorf("missing required fields at line %d", lineNumber)
        }

        if !strings.Contains(record.Email, "@") {
            return nil, fmt.Errorf("invalid email format at line %d", lineNumber)
        }

        records = append(records, record)
    }

    if len(records) == 0 {
        return nil, fmt.Errorf("no valid records found in file")
    }

    return records, nil
}

func ValidateRecords(records []DataRecord) []string {
    var errors []string
    emailSet := make(map[string]bool)

    for i, record := range records {
        if record.Active != "true" && record.Active != "false" {
            errors = append(errors, fmt.Sprintf("record %d: invalid active status '%s'", i+1, record.Active))
        }

        if emailSet[record.Email] {
            errors = append(errors, fmt.Sprintf("record %d: duplicate email '%s'", i+1, record.Email))
        }
        emailSet[record.Email] = true
    }

    return errors
}

func GenerateReport(records []DataRecord) {
    activeCount := 0
    for _, record := range records {
        if record.Active == "true" {
            activeCount++
        }
    }

    fmt.Printf("Total records processed: %d\n", len(records))
    fmt.Printf("Active records: %d\n", activeCount)
    fmt.Printf("Inactive records: %d\n", len(records)-activeCount)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: data_processor <csv_file_path>")
        os.Exit(1)
    }

    filePath := os.Args[1]
    records, err := ProcessCSVFile(filePath)
    if err != nil {
        fmt.Printf("Error processing file: %v\n", err)
        os.Exit(1)
    }

    if validationErrors := ValidateRecords(records); len(validationErrors) > 0 {
        fmt.Println("Validation errors found:")
        for _, errMsg := range validationErrors {
            fmt.Printf("  - %s\n", errMsg)
        }
        os.Exit(1)
    }

    GenerateReport(records)
    fmt.Println("Data processing completed successfully")
}