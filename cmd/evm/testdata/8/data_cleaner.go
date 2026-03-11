
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func removeDuplicates(inputFile, outputFile string) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	reader := csv.NewReader(in)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	var uniqueRecords [][]string

	for _, record := range records {
		if len(record) == 0 {
			continue
		}
		key := record[0]
		for i := 1; i < len(record); i++ {
			key += "," + record[i]
		}
		if !seen[key] {
			seen[key] = true
			uniqueRecords = append(uniqueRecords, record)
		}
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	return writer.WriteAll(uniqueRecords)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	err := removeDuplicates(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Duplicate removal completed successfully")
}
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

func DeduplicateEmails(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		email := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[email] {
			seen[email] = true
			record.Email = email
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmailFormat(email string) bool {
	if len(email) == 0 {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func CleanData(records []DataRecord) []DataRecord {
	deduped := DeduplicateEmails(records)
	var cleaned []DataRecord

	for _, record := range deduped {
		record.Valid = ValidateEmailFormat(record.Email)
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "user@example.com", false},
	}

	cleaned := CleanData(sampleData)

	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}
package main

import (
	"encoding/csv"
	"io"
	"strings"
)

func CleanCSVData(input io.Reader, output io.Writer) error {
	reader := csv.NewReader(input)
	writer := csv.NewWriter(output)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := make([]string, 0, len(record))
		hasData := false

		for _, field := range record {
			trimmed := strings.TrimSpace(field)
			cleanedRecord = append(cleanedRecord, trimmed)
			if trimmed != "" {
				hasData = true
			}
		}

		if hasData {
			if err := writer.Write(cleanedRecord); err != nil {
				return err
			}
		}
	}

	return nil
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
    ID      int
    Name    string
    Email   string
    Age     int
    Active  bool
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

    lineNum := 1
    for {
        lineNum++
        row, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Printf("line %d: skipped due to read error: %v\n", lineNum, err)
            continue
        }

        if len(row) != 5 {
            fmt.Printf("line %d: skipped due to column count mismatch\n", lineNum)
            continue
        }

        record, err := parseRecord(row)
        if err != nil {
            fmt.Printf("line %d: skipped due to validation error: %v\n", lineNum, err)
            continue
        }

        if !isValidRecord(record) {
            fmt.Printf("line %d: skipped due to invalid data\n", lineNum)
            continue
        }

        cleanRow := formatRecord(record)
        if err := writer.Write(cleanRow); err != nil {
            return fmt.Errorf("failed to write record: %w", err)
        }
    }

    return nil
}

func parseRecord(row []string) (Record, error) {
    var rec Record
    var err error

    if rec.ID, err = strconv.Atoi(strings.TrimSpace(row[0])); err != nil {
        return rec, fmt.Errorf("invalid ID: %w", err)
    }

    rec.Name = strings.TrimSpace(row[1])
    if rec.Name == "" {
        return rec, fmt.Errorf("empty name")
    }

    rec.Email = strings.TrimSpace(row[2])
    if !strings.Contains(rec.Email, "@") {
        return rec, fmt.Errorf("invalid email format")
    }

    if rec.Age, err = strconv.Atoi(strings.TrimSpace(row[3])); err != nil {
        return rec, fmt.Errorf("invalid age: %w", err)
    }

    rec.Active = strings.ToLower(strings.TrimSpace(row[4])) == "true"
    return rec, nil
}

func isValidRecord(r Record) bool {
    return r.ID > 0 &&
        len(r.Name) <= 100 &&
        len(r.Email) <= 255 &&
        r.Age >= 0 && r.Age <= 150
}

func formatRecord(r Record) []string {
    activeStr := "false"
    if r.Active {
        activeStr = "true"
    }
    return []string{
        strconv.Itoa(r.ID),
        r.Name,
        strings.ToLower(r.Email),
        strconv.Itoa(r.Age),
        activeStr,
    }
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
        os.Exit(1)
    }

    inputFile := os.Args[1]
    outputFile := os.Args[2]

    if err := cleanCSV(inputFile, outputFile); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Data cleaning completed successfully")
}