
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Record struct {
	ID    string
	Email string
	Phone string
	Valid bool
}

func deduplicateByHash(records []Record) []Record {
	seen := make(map[string]bool)
	var unique []Record

	for _, r := range records {
		hash := generateRecordHash(r)
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, r)
		}
	}
	return unique
}

func generateRecordHash(r Record) string {
	data := fmt.Sprintf("%s-%s-%s", r.ID, strings.ToLower(r.Email), r.Phone)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func validateRecords(records []Record) []Record {
	var valid []Record
	for _, r := range records {
		if isValidEmail(r.Email) && isValidPhone(r.Phone) {
			r.Valid = true
			valid = append(valid, r)
		}
	}
	return valid
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidPhone(phone string) bool {
	return len(phone) >= 10 && len(phone) <= 15
}

func processDataPipeline(records []Record) []Record {
	cleaned := deduplicateByHash(records)
	validated := validateRecords(cleaned)
	return validated
}

func main() {
	sampleData := []Record{
		{"001", "user@example.com", "1234567890", false},
		{"002", "user@example.com", "1234567890", false},
		{"003", "invalid-email", "0987654321", false},
		{"004", "another@test.org", "5551234567", false},
	}

	result := processDataPipeline(sampleData)
	fmt.Printf("Processed %d records, %d valid unique records found\n", len(sampleData), len(result))
	for _, r := range result {
		fmt.Printf("Valid record: ID=%s, Email=%s\n", r.ID, r.Email)
	}
}package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// Remove leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Remove any HTML/XML tags
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(trimmed, "")
	
	// Escape potentially dangerous characters
	re = regexp.MustCompile(`['"\\;]`)
	escaped := re.ReplaceAllStringFunc(cleaned, func(match string) string {
		return "\\" + match
	})
	
	// Limit length to prevent buffer overflow attacks
	maxLength := 1000
	if len(escaped) > maxLength {
		escaped = escaped[:maxLength]
	}
	
	return escaped
}