
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	processedRecords map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		processedRecords: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(records []string) []string {
	var uniqueRecords []string
	for _, record := range records {
		normalized := strings.ToLower(strings.TrimSpace(record))
		if !dc.processedRecords[normalized] {
			dc.processedRecords[normalized] = true
			uniqueRecords = append(uniqueRecords, record)
		}
	}
	return uniqueRecords
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && len(parts[1]) > 0 && strings.Contains(parts[1], ".")
}

func (dc *DataCleaner) CleanPhoneNumber(phone string) string {
	var cleaned strings.Builder
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			cleaned.WriteRune(ch)
		}
	}
	return cleaned.String()
}

func main() {
	cleaner := NewDataCleaner()
	
	records := []string{"user1@example.com", "User1@Example.com", "user2@test.org", "  user1@example.com  "}
	unique := cleaner.RemoveDuplicates(records)
	fmt.Println("Unique records:", unique)
	
	emails := []string{"test@domain.com", "invalid", "no@tld", "valid@address.co.uk"}
	for _, email := range emails {
		fmt.Printf("Email %s valid: %v\n", email, cleaner.ValidateEmail(email))
	}
	
	phoneNumbers := []string{"+1 (123) 456-7890", "123.456.7890", "123-456-7890"}
	for _, phone := range phoneNumbers {
		fmt.Printf("Cleaned phone: %s\n", cleaner.CleanPhoneNumber(phone))
	}
}package main

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
	input := []int{4, 2, 7, 2, 4, 9, 7}
	output := removeDuplicates(input)
	fmt.Println("Original:", input)
	fmt.Println("Deduplicated:", output)
}package main

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
}