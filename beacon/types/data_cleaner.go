
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