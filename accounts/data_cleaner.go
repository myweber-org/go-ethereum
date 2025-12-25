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

type DataCleaner struct {
    records []DataRecord
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        records: make([]DataRecord, 0),
    }
}

func (dc *DataCleaner) AddRecord(id int, email string) {
    record := DataRecord{
        ID:    id,
        Email: strings.ToLower(strings.TrimSpace(email)),
        Valid: strings.Contains(email, "@"),
    }
    dc.records = append(dc.records, record)
}

func (dc *DataCleaner) RemoveDuplicates() []DataRecord {
    seen := make(map[string]bool)
    unique := make([]DataRecord, 0)

    for _, record := range dc.records {
        if !seen[record.Email] {
            seen[record.Email] = true
            unique = append(unique, record)
        }
    }

    dc.records = unique
    return unique
}

func (dc *DataCleaner) GetValidRecords() []DataRecord {
    valid := make([]DataRecord, 0)
    for _, record := range dc.records {
        if record.Valid {
            valid = append(valid, record)
        }
    }
    return valid
}

func (dc *DataCleaner) PrintSummary() {
    fmt.Printf("Total records: %d\n", len(dc.records))
    fmt.Printf("Valid records: %d\n", len(dc.GetValidRecords()))
}

func main() {
    cleaner := NewDataCleaner()
    
    cleaner.AddRecord(1, "user@example.com")
    cleaner.AddRecord(2, "user@example.com")
    cleaner.AddRecord(3, "admin@test.org")
    cleaner.AddRecord(4, "invalid-email")
    cleaner.AddRecord(5, "  TEST@DOMAIN.COM  ")
    
    fmt.Println("Before deduplication:")
    cleaner.PrintSummary()
    
    unique := cleaner.RemoveDuplicates()
    fmt.Println("\nAfter deduplication:")
    fmt.Printf("Unique records: %d\n", len(unique))
    
    fmt.Println("\nValid records:")
    for _, record := range cleaner.GetValidRecords() {
        fmt.Printf("ID: %d, Email: %s\n", record.ID, record.Email)
    }
}