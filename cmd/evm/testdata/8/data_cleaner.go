
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
}