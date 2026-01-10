
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

func DeduplicateEmails(emails []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, email := range emails {
        email = strings.ToLower(strings.TrimSpace(email))
        if !seen[email] {
            seen[email] = true
            result = append(result, email)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    if len(email) < 3 || !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
        return false
    }
    return true
}

func CleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        cleanEmail := strings.ToLower(strings.TrimSpace(record.Email))
        if ValidateEmail(cleanEmail) && !emailSet[cleanEmail] {
            emailSet[cleanEmail] = true
            record.Email = cleanEmail
            record.Valid = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", false},
        {2, "USER@example.com", false},
        {3, "invalid-email", false},
        {4, "another@test.org", false},
        {5, "user@example.com", false},
    }
    
    cleaned := CleanData(sampleData)
    fmt.Printf("Original: %d records\n", len(sampleData))
    fmt.Printf("Cleaned: %d records\n", len(cleaned))
    
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", record.ID, record.Email, record.Valid)
    }
}