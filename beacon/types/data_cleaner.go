
package main

import (
    "fmt"
    "strings"
)

// DataCleaner provides methods for cleaning datasets
type DataCleaner struct{}

// RemoveDuplicates removes duplicate strings from a slice
func (dc DataCleaner) RemoveDuplicates(input []string) []string {
    encountered := map[string]bool{}
    result := []string{}

    for _, value := range input {
        if !encountered[value] {
            encountered[value] = true
            result = append(result, value)
        }
    }
    return result
}

// ValidateEmail checks if a string is a valid email format
func (dc DataCleaner) ValidateEmail(email string) bool {
    if len(email) < 3 || len(email) > 254 {
        return false
    }
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// NormalizeWhitespace replaces multiple spaces with single spaces
func (dc DataCleaner) NormalizeWhitespace(text string) string {
    return strings.Join(strings.Fields(text), " ")
}

func main() {
    cleaner := DataCleaner{}

    // Example usage
    data := []string{"apple", "banana", "apple", "cherry", "banana"}
    unique := cleaner.RemoveDuplicates(data)
    fmt.Printf("Original: %v\n", data)
    fmt.Printf("Deduplicated: %v\n", unique)

    emails := []string{"test@example.com", "invalid-email", "user@domain.co.uk"}
    for _, email := range emails {
        fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
    }

    text := "  This   has   extra   spaces  "
    normalized := cleaner.NormalizeWhitespace(text)
    fmt.Printf("Original: '%s'\n", text)
    fmt.Printf("Normalized: '%s'\n", normalized)
}
package main

import (
	"fmt"
	"strings"
)

type DataRecord struct {
	ID    int
	Email string
	Name  string
}

type DataCleaner struct {
	records []DataRecord
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		records: make([]DataRecord, 0),
	}
}

func (dc *DataCleaner) AddRecord(record DataRecord) {
	dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
	seen := make(map[string]DataRecord)
	result := make([]DataRecord, 0)

	for _, record := range dc.records {
		key := fmt.Sprintf("%s|%s", strings.ToLower(record.Email), strings.ToLower(record.Name))
		if _, exists := seen[key]; !exists {
			seen[key] = record
			result = append(result, record)
		}
	}

	dc.records = result
	return result
}

func (dc *DataCleaner) ValidateEmails() (valid []DataRecord, invalid []DataRecord) {
	valid = make([]DataRecord, 0)
	invalid = make([]DataRecord, 0)

	for _, record := range dc.records {
		if strings.Contains(record.Email, "@") && strings.Contains(record.Email, ".") {
			valid = append(valid, record)
		} else {
			invalid = append(invalid, record)
		}
	}

	return valid, invalid
}

func (dc *DataCleaner) GetRecordCount() int {
	return len(dc.records)
}

func main() {
	cleaner := NewDataCleaner()

	cleaner.AddRecord(DataRecord{ID: 1, Email: "user@example.com", Name: "John Doe"})
	cleaner.AddRecord(DataRecord{ID: 2, Email: "user@example.com", Name: "John Doe"})
	cleaner.AddRecord(DataRecord{ID: 3, Email: "jane@test.org", Name: "Jane Smith"})
	cleaner.AddRecord(DataRecord{ID: 4, Email: "invalid-email", Name: "Bad Data"})

	fmt.Printf("Initial records: %d\n", cleaner.GetRecordCount())

	cleaner.RemoveDuplicates()
	fmt.Printf("After deduplication: %d\n", cleaner.GetRecordCount())

	valid, invalid := cleaner.ValidateEmails()
	fmt.Printf("Valid emails: %d, Invalid emails: %d\n", len(valid), len(invalid))

	for _, record := range valid {
		fmt.Printf("Valid: %s <%s>\n", record.Name, record.Email)
	}
}