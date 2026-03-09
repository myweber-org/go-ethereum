package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Name  string
	Email string
	Valid bool
}

func DeduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := fmt.Sprintf("%s|%s", record.Name, record.Email)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func MarkValidRecords(records []DataRecord) []DataRecord {
	for i := range records {
		records[i].Valid = ValidateEmail(records[i].Email)
	}
	return records
}

func ProcessData(input []DataRecord) []DataRecord {
	deduped := DeduplicateRecords(input)
	validated := MarkValidRecords(deduped)
	return validated
}

func main() {
	sampleData := []DataRecord{
		{1, "John Doe", "john@example.com", false},
		{2, "Jane Smith", "jane@example.com", false},
		{3, "John Doe", "john@example.com", false},
		{4, "Bob Wilson", "invalid-email", false},
	}

	processed := ProcessData(sampleData)

	for _, record := range processed {
		status := "INVALID"
		if record.Valid {
			status = "VALID"
		}
		fmt.Printf("ID: %d, Name: %s, Status: %s\n", record.ID, record.Name, status)
	}
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

func parseCSV(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records []Record
	lineNum := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lineNum++

		if lineNum == 1 {
			continue
		}

		if len(line) != 4 {
			continue
		}

		id, err := strconv.Atoi(strings.TrimSpace(line[0]))
		if err != nil {
			continue
		}

		name := strings.TrimSpace(line[1])
		email := strings.TrimSpace(line[2])
		score, err := strconv.ParseFloat(strings.TrimSpace(line[3]), 64)
		if err != nil {
			continue
		}

		if !validateEmail(email) {
			continue
		}

		records = append(records, Record{
			ID:    id,
			Name:  name,
			Email: email,
			Score: score,
		})
	}

	return records, nil
}

func validateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func calculateStats(records []Record) (float64, float64, int) {
	if len(records) == 0 {
		return 0, 0, 0
	}

	var sum float64
	var maxScore float64
	validCount := 0

	for _, record := range records {
		if record.Score >= 0 && record.Score <= 100 {
			sum += record.Score
			if record.Score > maxScore {
				maxScore = record.Score
			}
			validCount++
		}
	}

	average := sum / float64(validCount)
	return average, maxScore, validCount
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: data_cleaner <csv_file>")
		return
	}

	records, err := parseCSV(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		return
	}

	fmt.Printf("Successfully parsed %d records\n", len(records))

	avgScore, maxScore, validCount := calculateStats(records)
	fmt.Printf("Valid records: %d\n", validCount)
	fmt.Printf("Average score: %.2f\n", avgScore)
	fmt.Printf("Maximum score: %.2f\n", maxScore)
}