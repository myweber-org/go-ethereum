package main

import (
	"regexp"
	"strings"
)

func NormalizeEmail(email string) (string, bool) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil || !matched {
		return "", false
	}
	return email, true
}

func ValidateUsername(username string) bool {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return false
	}
	return matched
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`<[^>]*>`)
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

type DataRecord struct {
	ID    string
	Name  string
	Email string
	Valid bool
}

func ProcessCSVFile(filename string) ([]DataRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records := []DataRecord{}
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

		if len(row) < 3 {
			continue
		}

		record := DataRecord{
			ID:    strings.TrimSpace(row[0]),
			Name:  strings.TrimSpace(row[1]),
			Email: strings.TrimSpace(row[2]),
			Valid: validateRecord(strings.TrimSpace(row[0]), strings.TrimSpace(row[2])),
		}

		if record.Valid {
			records = append(records, record)
		}
	}

	return records, nil
}

func validateRecord(id, email string) bool {
	if id == "" || email == "" {
		return false
	}
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func GenerateReport(records []DataRecord) {
	fmt.Printf("Total valid records: %d\n", len(records))
	for _, record := range records {
		fmt.Printf("ID: %s, Name: %s, Email: %s\n", record.ID, record.Name, record.Email)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_processor <csv_file>")
		return
	}

	records, err := ProcessCSVFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		return
	}

	GenerateReport(records)
}