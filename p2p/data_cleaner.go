
package main

import (
    "fmt"
    "strings"
)

func RemoveDuplicates(slice []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    for _, item := range slice {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    return result
}

func ValidateEmail(email string) bool {
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

func CleanData(data []string) []string {
    cleaned := RemoveDuplicates(data)
    validEmails := []string{}
    for _, item := range cleaned {
        if ValidateEmail(item) {
            validEmails = append(validEmails, item)
        }
    }
    return validEmails
}

func main() {
    sampleData := []string{
        "user@example.com",
        "admin@test.org",
        "user@example.com",
        "invalid-email",
        "another@domain.com",
        "",
        "missing@domain",
    }

    result := CleanData(sampleData)
    fmt.Println("Cleaned data:", result)
    fmt.Println("Original count:", len(sampleData))
    fmt.Println("Cleaned count:", len(result))
}
package main

import "fmt"

func removeDuplicates(nums []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

func main() {
	data := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	cleaned := removeDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}