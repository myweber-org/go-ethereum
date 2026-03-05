package main

import (
	"fmt"
	"strings"
)

func CleanData(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range input {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}
	return result
}

func main() {
	data := []string{"Apple", " apple ", "banana", "Banana", "Cherry"}
	cleaned := CleanData(data)
	fmt.Println("Cleaned data:", cleaned)
}package main

import "fmt"

func removeDuplicates(input []int) []int {
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
	slice := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
	uniqueSlice := removeDuplicates(slice)
	fmt.Println("Original:", slice)
	fmt.Println("Unique:", uniqueSlice)
}
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

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func CleanData(data []string) []string {
	normalized := make([]string, len(data))
	for i, item := range data {
		normalized[i] = NormalizeString(item)
	}
	return RemoveDuplicates(normalized)
}

func main() {
	sampleData := []string{"  Apple", "banana  ", "Apple", "  BANANA  ", "Cherry"}
	cleaned := CleanData(sampleData)
	fmt.Println("Cleaned Data:", cleaned)
}