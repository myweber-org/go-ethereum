
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
        email = strings.ToLower(strings.TrimSpace(email))
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

func cleanData(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        record.Email = strings.ToLower(strings.TrimSpace(record.Email))
        
        if validateEmail(record.Email) && !emailSet[record.Email] {
            record.Valid = true
            emailSet[record.Email] = true
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
        {4, "user@example.com", false},
        {5, "new.user@domain.org", false},
    }
    
    cleaned := cleanData(sampleData)
    
    fmt.Printf("Original records: %d\n", len(sampleData))
    fmt.Printf("Cleaned records: %d\n", len(cleaned))
    
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Valid: %v\n", 
            record.ID, record.Email, record.Valid)
    }
    
    emails := []string{"test@mail.com", "TEST@mail.com", "test@mail.com", "another@test.org"}
    uniqueEmails := deduplicateEmails(emails)
    fmt.Println("\nUnique emails:", uniqueEmails)
}