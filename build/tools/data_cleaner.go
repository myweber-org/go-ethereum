package utils

import (
	"regexp"
	"strings"
)

func SanitizeInput(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)
	
	// Replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(trimmed, " ")
	
	return cleaned
}
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
}

type DataCleaner struct {
	seenHashes map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seenHashes: make(map[string]bool),
	}
}

func (dc *DataCleaner) NormalizeEmail(email string) string {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(email)), "@")
	if len(parts) != 2 {
		return ""
	}
	local := strings.Split(parts[0], "+")[0]
	local = strings.ReplaceAll(local, ".", "")
	return local + "@" + parts[1]
}

func (dc *DataCleaner) GenerateHash(record Record) string {
	normalizedEmail := dc.NormalizeEmail(record.Email)
	normalizedPhone := strings.ReplaceAll(strings.TrimSpace(record.Phone), " ", "")

	hashInput := fmt.Sprintf("%s|%s|%s",
		record.ID,
		normalizedEmail,
		normalizedPhone,
	)

	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])
}

func (dc *DataCleaner) IsDuplicate(record Record) bool {
	hash := dc.GenerateHash(record)
	if dc.seenHashes[hash] {
		return true
	}
	dc.seenHashes[hash] = true
	return false
}

func (dc *DataCleaner) ValidateRecord(record Record) bool {
	if strings.TrimSpace(record.ID) == "" {
		return false
	}
	if dc.NormalizeEmail(record.Email) == "" {
		return false
	}
	if len(strings.ReplaceAll(record.Phone, " ", "")) < 8 {
		return false
	}
	return true
}

func (dc *DataCleaner) ProcessRecords(records []Record) []Record {
	var cleaned []Record
	for _, record := range records {
		if !dc.ValidateRecord(record) {
			continue
		}
		if dc.IsDuplicate(record) {
			continue
		}
		cleaned = append(cleaned, record)
	}
	return cleaned
}

func main() {
	cleaner := NewDataCleaner()

	records := []Record{
		{"001", "user@example.com", "12345678"},
		{"002", "USER@example.com", "12345678"},
		{"003", "user+tag@example.com", "12345678"},
		{"004", "u.s.e.r@example.com", "12345678"},
		{"005", "invalid-email", "12345678"},
		{"006", "another@example.com", ""},
		{"007", "unique@example.com", "87654321"},
	}

	cleaned := cleaner.ProcessRecords(records)
	fmt.Printf("Original: %d, Cleaned: %d\n", len(records), len(cleaned))
}