
package main

import (
	"errors"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Valid bool
}

func RemoveDuplicates(records []DataRecord) []DataRecord {
	seen := make(map[int]bool)
	var unique []DataRecord

	for _, record := range records {
		if !seen[record.ID] {
			seen[record.ID] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

func FilterValidRecords(records []DataRecord) []DataRecord {
	var valid []DataRecord
	for _, record := range records {
		if record.Valid && ValidateEmail(record.Email) == nil {
			valid = append(valid, record)
		}
	}
	return valid
}

func main() {
	records := []DataRecord{
		{1, "user@example.com", true},
		{2, "invalid-email", false},
		{1, "duplicate@example.com", true},
		{3, "another@test.org", true},
	}

	fmt.Println("Original records:", len(records))
	unique := RemoveDuplicates(records)
	fmt.Println("After deduplication:", len(unique))
	valid := FilterValidRecords(unique)
	fmt.Println("Valid records:", len(valid))
}