
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
    result := []DataRecord{}
    for _, record := range records {
        normalizedEmail := strings.ToLower(strings.TrimSpace(record.Email))
        if !seen[normalizedEmail] {
            seen[normalizedEmail] = true
            result = append(result, record)
        }
    }
    return result
}

func ValidateEmails(records []DataRecord) []DataRecord {
    for i := range records {
        email := records[i].Email
        records[i].Valid = strings.Contains(email, "@") && strings.Contains(email, ".")
    }
    return records
}

func CleanData(records []DataRecord) []DataRecord {
    deduped := RemoveDuplicates(records)
    validated := ValidateEmails(deduped)
    return validated
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "USER@example.com", false},
        {3, "invalid-email", false},
        {4, "test@domain.org", false},
    }

    cleaned := CleanData(sampleData)
    fmt.Printf("Processed %d records\n", len(cleaned))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
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

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord

	for _, record := range records {
		key := strings.ToLower(strings.TrimSpace(record.Email))
		if !seen[key] {
			seen[key] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return true
}

func cleanData(records []DataRecord) []DataRecord {
	var cleaned []DataRecord
	unique := deduplicateRecords(records)

	for _, record := range unique {
		record.Valid = validateEmail(record.Email)
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	sampleData := []DataRecord{
		{1, "user@example.com", false},
		{2, "user@example.com", false},
		{3, "invalid-email", false},
		{4, "test@domain.org", false},
		{5, "ANOTHER@EXAMPLE.COM", false},
	}

	cleaned := cleanData(sampleData)
	for _, record := range cleaned {
		fmt.Printf("ID: %d, Email: %s, Valid: %t\n", record.ID, record.Email, record.Valid)
	}
}