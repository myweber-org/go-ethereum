
package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package main

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
    return strings.Contains(parts[1], ".")
}

func CleanRecords(records []DataRecord) []DataRecord {
    emailSet := make(map[string]bool)
    cleaned := []DataRecord{}
    
    for _, record := range records {
        record.Email = strings.ToLower(strings.TrimSpace(record.Email))
        record.Valid = ValidateEmail(record.Email)
        
        if record.Valid && !emailSet[record.Email] {
            emailSet[record.Email] = true
            cleaned = append(cleaned, record)
        }
    }
    return cleaned
}

func main() {
    emails := []string{
        "test@example.com",
        "TEST@example.com",
        "invalid-email",
        "another@test.org",
        "test@example.com",
    }
    
    uniqueEmails := DeduplicateEmails(emails)
    fmt.Printf("Unique emails: %v\n", uniqueEmails)
    
    records := []DataRecord{
        {1, "user@domain.com", false},
        {2, "DUPLICATE@domain.com", false},
        {3, "duplicate@domain.com", false},
        {4, "bad-email", false},
    }
    
    cleaned := CleanRecords(records)
    fmt.Printf("Cleaned records: %+v\n", cleaned)
}