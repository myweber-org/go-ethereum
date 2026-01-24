
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var uniqueRecords []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			uniqueRecords = append(uniqueRecords, record)
		}
	}
	return uniqueRecords
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) CleanPhoneNumber(phone string) string {
	var cleaned strings.Builder
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			cleaned.WriteRune(ch)
		}
	}
	return cleaned.String()
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{"user1@example.com", "User1@Example.com", "user2@test.org", "  user1@example.com  "}
	unique := cleaner.RemoveDuplicates(records)
	fmt.Println("Unique records:", unique)
	
	emails := []string{"test@domain.com", "invalid", "no@tld", "valid@address.co.uk"}
	for _, email := range emails {
		fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
	}
	
	phoneNumbers := []string{"+1 (123) 456-7890", "123.456.7890", "123-456-7890"}
	for _, phone := range phoneNumbers {
		fmt.Printf("Cleaned phone: %s\n", cleaner.CleanPhoneNumber(phone))
	}
}