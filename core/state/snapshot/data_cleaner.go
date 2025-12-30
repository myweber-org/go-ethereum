
package main

import (
    "fmt"
    "strings"
)

// DataCleaner provides methods for cleaning string data
type DataCleaner struct{}

// Deduplicate removes duplicate entries from a slice of strings
func (dc *DataCleaner) Deduplicate(items []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, item := range items {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    return result
}

// ValidateEmail checks if a string is a valid email format
func (dc *DataCleaner) ValidateEmail(email string) bool {
    if len(email) < 3 || len(email) > 254 {
        return false
    }
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    if len(parts[0]) == 0 || len(parts[1]) == 0 {
        return false
    }
    return true
}

// TrimSpaces removes leading and trailing whitespace from all strings
func (dc *DataCleaner) TrimSpaces(items []string) []string {
    trimmed := make([]string, len(items))
    for i, item := range items {
        trimmed[i] = strings.TrimSpace(item)
    }
    return trimmed
}

func main() {
    cleaner := &DataCleaner{}
    
    sampleData := []string{"  alice@example.com", "bob@test.org", "  alice@example.com", "invalid-email", "  charlie@demo.net  "}
    
    fmt.Println("Original data:", sampleData)
    
    trimmed := cleaner.TrimSpaces(sampleData)
    fmt.Println("After trimming:", trimmed)
    
    deduplicated := cleaner.Deduplicate(trimmed)
    fmt.Println("After deduplication:", deduplicated)
    
    fmt.Println("\nEmail validation results:")
    for _, email := range deduplicated {
        isValid := cleaner.ValidateEmail(email)
        fmt.Printf("%s: %v\n", email, isValid)
    }
}