package main

import "fmt"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	numbers := []int{1, 2, 2, 3, 4, 4, 5}
	uniqueNumbers := RemoveDuplicates(numbers)
	fmt.Println("Original:", numbers)
	fmt.Println("Unique:", uniqueNumbers)

	strings := []string{"apple", "banana", "apple", "orange"}
	uniqueStrings := RemoveDuplicates(strings)
	fmt.Println("Original:", strings)
	fmt.Println("Unique:", uniqueStrings)
}package main

import (
    "fmt"
    "strings"
)

type DataCleaner struct {
    seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        seen: make(map[string]bool),
    }
}

func (dc *DataCleaner) Clean(input string) string {
    trimmed := strings.TrimSpace(input)
    if trimmed == "" {
        return ""
    }
    lower := strings.ToLower(trimmed)
    if dc.seen[lower] {
        return ""
    }
    dc.seen[lower] = true
    return trimmed
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    return len(parts[0]) > 0 && len(parts[1]) > 0
}

func main() {
    cleaner := NewDataCleaner()
    samples := []string{
        "  User@Example.com  ",
        "user@example.com",
        "invalid-email",
        "",
        "  another@test.org  ",
    }
    for _, sample := range samples {
        cleaned := cleaner.Clean(sample)
        if cleaned != "" {
            valid := cleaner.ValidateEmail(cleaned)
            fmt.Printf("Original: %q -> Cleaned: %q (Valid: %v)\n", sample, cleaned, valid)
        }
    }
}