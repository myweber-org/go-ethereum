
package main

import (
    "fmt"
    "strings"
)

type DataRecord struct {
    ID    int
    Email string
    Phone string
}

func deduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[int]bool)
    var result []DataRecord
    for _, record := range records {
        if !seen[record.ID] {
            seen[record.ID] = true
            result = append(result, record)
        }
    }
    return result
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func cleanPhoneNumber(phone string) string {
    var cleaned strings.Builder
    for _, ch := range phone {
        if ch >= '0' && ch <= '9' {
            cleaned.WriteRune(ch)
        }
    }
    return cleaned.String()
}

func processRecords(records []DataRecord) []DataRecord {
    var validRecords []DataRecord
    for _, record := range deduplicateRecords(records) {
        if validateEmail(record.Email) {
            record.Phone = cleanPhoneNumber(record.Phone)
            validRecords = append(validRecords, record)
        }
    }
    return validRecords
}

func main() {
    sampleData := []DataRecord{
        {1, "test@example.com", "(123) 456-7890"},
        {2, "invalid-email", "555-1234"},
        {1, "test@example.com", "1234567890"},
        {3, "user@domain.org", "+1-800-555-0199"},
    }

    cleaned := processRecords(sampleData)
    fmt.Printf("Processed %d records\n", len(cleaned))
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", record.ID, record.Email, record.Phone)
    }
}