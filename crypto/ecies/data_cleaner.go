package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type DataRecord struct {
	ID   string
	Name string
	Email string
}

func generateHash(record DataRecord) string {
	data := fmt.Sprintf("%s|%s|%s", record.ID, record.Name, record.Email)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func deduplicateRecords(records []DataRecord) []DataRecord {
	seen := make(map[string]bool)
	var unique []DataRecord
	
	for _, record := range records {
		hash := generateHash(record)
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, record)
		}
	}
	return unique
}

func validateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0
}

func cleanData(records []DataRecord) []DataRecord {
	var valid []DataRecord
	
	for _, record := range records {
		if validateEmail(record.Email) {
			valid = append(valid, record)
		}
	}
	
	return deduplicateRecords(valid)
}

func main() {
	sampleData := []DataRecord{
		{"1", "John Doe", "john@example.com"},
		{"2", "Jane Smith", "jane@example.com"},
		{"3", "John Doe", "john@example.com"},
		{"4", "Bob Wilson", "invalid-email"},
		{"5", "Alice Brown", "alice@example.com"},
		{"6", "Bob Wilson", "bob@example.com"},
		{"7", "Bob Wilson", "bob@example.com"},
	}
	
	cleaned := cleanData(sampleData)
	
	fmt.Printf("Original records: %d\n", len(sampleData))
	fmt.Printf("Cleaned records: %d\n", len(cleaned))
	
	for i, record := range cleaned {
		fmt.Printf("%d: %s - %s\n", i+1, record.Name, record.Email)
	}
}