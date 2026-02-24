
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
	var unique []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) CleanPhoneNumber(phone string) string {
	var builder strings.Builder
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{
		"john@example.com",
		"Jane@Example.COM",
		"john@example.com",
		"invalid-email",
		"  ALICE@DOMAIN.COM  ",
	}
	
	fmt.Println("Original records:", records)
	unique := cleaner.RemoveDuplicates(records)
	fmt.Println("Deduplicated records:", unique)
	
	testEmails := []string{"test@domain.com", "invalid", "user@com", "@domain.com"}
	for _, email := range testEmails {
		fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
	}
	
	phone := "+1 (234) 567-8900"
	cleanedPhone := cleaner.CleanPhoneNumber(phone)
	fmt.Printf("Phone %s cleaned: %s\n", phone, cleanedPhone)
}