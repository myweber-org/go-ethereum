package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func CleanCSVRow(row []string) []string {
	cleaned := make([]string, len(row))
	for i, field := range row {
		trimmed := strings.TrimSpace(field)
		normalized := strings.ToLower(trimmed)
		cleaned[i] = normalized
	}
	return cleaned
}

func ProcessCSV(reader io.Reader, writer io.Writer) error {
	csvReader := csv.NewReader(reader)
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		cleanedRecord := CleanCSVRow(record)
		if err := csvWriter.Write(cleanedRecord); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	sampleInput := "Name,Age,Email\nJohn Doe,25,JOHN@example.com\n Jane Smith,30,jane@test.org "
	reader := strings.NewReader(sampleInput)
	var output strings.Builder

	err := ProcessCSV(reader, &output)
	if err != nil {
		fmt.Printf("Error processing CSV: %v\n", err)
		return
	}

	fmt.Println("Cleaned CSV output:")
	fmt.Println(output.String())
}
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
    if len(nums) == 0 {
        return nums
    }
    
    seen := make(map[int]bool)
    result := make([]int, 0, len(nums))
    
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    
    return result
}

func main() {
    testData := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := RemoveDuplicates(testData)
    fmt.Printf("Original: %v\n", testData)
    fmt.Printf("Cleaned: %v\n", cleaned)
}package main

import (
	"errors"
	"fmt"
	"strings"
)

type Record struct {
	ID    int
	Email string
	Valid bool
}

func DeduplicateEmails(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record
	for _, r := range records {
		email := strings.ToLower(strings.TrimSpace(r.Email))
		if !seen[email] {
			seen[email] = true
			r.Email = email
			unique = append(unique, r)
		}
	}
	return unique
}

func ValidateEmail(email string) error {
	if len(email) == 0 {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("email must contain @ symbol")
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return errors.New("invalid email format")
	}
	return nil
}

func CleanRecords(records []Record) ([]Record, error) {
	cleaned := DeduplicateEmails(records)
	for i := range cleaned {
		if err := ValidateEmail(cleaned[i].Email); err != nil {
			cleaned[i].Valid = false
			return cleaned, fmt.Errorf("record %d invalid: %w", cleaned[i].ID, err)
		}
		cleaned[i].Valid = true
	}
	return cleaned, nil
}

func main() {
	records := []Record{
		{1, "user@example.com", false},
		{2, "USER@example.com", false},
		{3, "invalid-email", false},
	}

	cleaned, err := CleanRecords(records)
	if err != nil {
		fmt.Printf("Cleaning error: %v\n", err)
	}
	fmt.Printf("Processed %d records\n", len(cleaned))
}