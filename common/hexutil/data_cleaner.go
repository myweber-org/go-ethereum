
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

func RemoveDuplicates(records []DataRecord) []DataRecord {
    seen := make(map[int]bool)
    result := []DataRecord{}
    
    for _, record := range records {
        if !seen[record.ID] {
            seen[record.ID] = true
            result = append(result, record)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
    if len(email) == 0 {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FormatPhoneNumber(phone string) string {
    cleaned := strings.ReplaceAll(phone, " ", "")
    cleaned = strings.ReplaceAll(cleaned, "-", "")
    cleaned = strings.ReplaceAll(cleaned, "(", "")
    cleaned = strings.ReplaceAll(cleaned, ")", "")
    
    if len(cleaned) == 10 {
        return fmt.Sprintf("(%s) %s-%s", 
            cleaned[0:3], cleaned[3:6], cleaned[6:10])
    }
    return phone
}

func CleanData(records []DataRecord) []DataRecord {
    uniqueRecords := RemoveDuplicates(records)
    
    for i := range uniqueRecords {
        if !ValidateEmail(uniqueRecords[i].Email) {
            uniqueRecords[i].Email = "invalid@example.com"
        }
        uniqueRecords[i].Phone = FormatPhoneNumber(uniqueRecords[i].Phone)
    }
    
    return uniqueRecords
}

func main() {
    sampleData := []DataRecord{
        {1, "user@domain.com", "1234567890"},
        {2, "invalid-email", "555-123-4567"},
        {1, "user@domain.com", "1234567890"},
        {3, "another@test.org", "(987)654-3210"},
    }
    
    cleaned := CleanData(sampleData)
    
    for _, record := range cleaned {
        fmt.Printf("ID: %d, Email: %s, Phone: %s\n", 
            record.ID, record.Email, record.Phone)
    }
}