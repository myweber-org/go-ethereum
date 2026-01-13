
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
}
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

func DeduplicateRecords(records []DataRecord) []DataRecord {
    seen := make(map[string]bool)
    unique := []DataRecord{}
    
    for _, record := range records {
        key := fmt.Sprintf("%s|%s", record.Email, record.Phone)
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
    if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
        return false
    }
    return true
}

func CleanPhoneNumber(phone string) string {
    var builder strings.Builder
    for _, ch := range phone {
        if ch >= '0' && ch <= '9' {
            builder.WriteRune(ch)
        }
    }
    return builder.String()
}

func ProcessRecords(records []DataRecord) []DataRecord {
    cleaned := []DataRecord{}
    
    for _, record := range records {
        if !ValidateEmail(record.Email) {
            continue
        }
        
        cleanedPhone := CleanPhoneNumber(record.Phone)
        if len(cleanedPhone) < 10 {
            continue
        }
        
        cleaned = append(cleaned, DataRecord{
            ID:    record.ID,
            Email: strings.ToLower(strings.TrimSpace(record.Email)),
            Phone: cleanedPhone,
        })
    }
    
    return DeduplicateRecords(cleaned)
}

func main() {
    sampleData := []DataRecord{
        {1, "user@example.com", "(123) 456-7890"},
        {2, "user@example.com", "123-456-7890"},
        {3, "invalid-email", "555-1234"},
        {4, "another@test.org", "987.654.3210"},
        {5, "ANOTHER@TEST.ORG", "9876543210"},
    }
    
    result := ProcessRecords(sampleData)
    fmt.Printf("Processed %d records, kept %d valid unique records\n", 
        len(sampleData), len(result))
    
    for _, record := range result {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", 
            record.ID, record.Email, record.Phone)
    }
}