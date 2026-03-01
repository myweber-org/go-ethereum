
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

func deduplicateEmails(emails []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, email := range emails {
        if !seen[email] {
            seen[email] = true
            result = append(result, email)
        }
    }
    return result
}

func validateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func processRecords(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    var cleaned []DataRecord

    for _, record := range records {
        if !emailSet[record.Email] && validateEmail(record.Email) {
            emailSet[record.Email] = true
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    records := []DataRecord{
        {1, "user@example.com", false},
        {2, "invalid-email", false},
        {3, "user@example.com", false},
        {4, "another@test.org", false},
    }

    cleaned := processRecords(records)
    fmt.Printf("Processed %d records, %d valid unique records found\n", len(records), len(cleaned))
    for _, r := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", r.ID, r.Email, r.Valid)
    }
}