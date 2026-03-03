
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmails(records []DataRecord) []DataRecord {
	var validated []DataRecord
	for _, record := range records {
		record.Valid = strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".")
		validated = append(validated, record)
	}
	return validated
}

func PrintRecords(records []DataRecord) {
	for _, record := range records {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Email: %s, Status: %s\n", record.ID, record.Email, status)
	}
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", false},
		{2, "test@domain.org", false},
		{3, "user@example.com", false},
		{4, "invalid-email", false},
		{5, "another@test.com", false},
		{6, "test@domain.org", false},
	}

	fmt.Println("Original records:")
	PrintRecords(records)

	uniqueRecords := RemoveDuplicates(records)
	fmt.Println("\nAfter deduplication:")
	PrintRecords(uniqueRecords)

	validatedRecords := ValidateEmails(uniqueRecords)
	fmt.Println("\nAfter validation:")
	PrintRecords(validatedRecords)
}package utils

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}package main

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

    headers, err := reader.Read()
    if err != nil {
        return fmt.Errorf("failed to read headers: %w", err)
    }

    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("failed to write headers: %w", err)
    }

    for {
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            continue
        }

        cleaned := processRow(row)
        if cleaned != nil {
            writer.Write([]string{
                strconv.Itoa(cleaned.ID),
                cleaned.Name,
                cleaned.Email,
                fmt.Sprintf("%.2f", cleaned.Score),
            })
        }
    }

    return nil
}

func processRow(row []string) *Record {
    if len(row) != 4 {
        return nil
    }

    id, err := strconv.Atoi(strings.TrimSpace(row[0]))
    if err != nil || id <= 0 {
        return nil
    }

    name := strings.TrimSpace(row[1])
    if name == "" || len(name) > 100 {
        return nil
    }

    email := strings.TrimSpace(row[2])
    if !strings.Contains(email, "@") || strings.Contains(email, " ") {
        return nil
    }

    score, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
    if err != nil || score < 0 || score > 100 {
        return nil
    }

    return &Record{
        ID:    id,
        Name:  name,
        Email: email,
        Score: score,
    }
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    if err := cleanCSV(os.Args[1], os.Args[2]); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}