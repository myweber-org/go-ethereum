package main

import "fmt"

func RemoveDuplicates(input []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, value := range input {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5}
	cleaned := RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}package main

import (
    "fmt"
    "strings"
)

type DataCleaner struct {
    duplicates map[string]int
}

func NewDataCleaner() *DataCleaner {
    return &DataCleaner{
        duplicates: make(map[string]int),
    }
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, item := range items {
        trimmed := strings.TrimSpace(item)
        if trimmed == "" {
            continue
        }
        if !seen[trimmed] {
            seen[trimmed] = true
            dc.duplicates[trimmed]++
            result = append(result, trimmed)
        } else {
            dc.duplicates[trimmed]++
        }
    }
    return result
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
    if !strings.Contains(email, "@") {
        return false
    }
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    if parts[0] == "" || parts[1] == "" {
        return false
    }
    return true
}

func (dc *DataCleaner) GetDuplicateStats() map[string]int {
    stats := make(map[string]int)
    for key, count := range dc.duplicates {
        if count > 1 {
            stats[key] = count
        }
    }
    return stats
}

func main() {
    cleaner := NewDataCleaner()
    
    sampleData := []string{
        "user@example.com",
        "user@example.com",
        "  user@example.com  ",
        "invalid-email",
        "another@test.org",
        "",
        "another@test.org",
    }
    
    cleaned := cleaner.RemoveDuplicates(sampleData)
    fmt.Println("Cleaned data:", cleaned)
    
    for _, email := range cleaned {
        if cleaner.ValidateEmail(email) {
            fmt.Printf("%s is valid\n", email)
        } else {
            fmt.Printf("%s is invalid\n", email)
        }
    }
    
    stats := cleaner.GetDuplicateStats()
    fmt.Println("Duplicate statistics:", stats)
}