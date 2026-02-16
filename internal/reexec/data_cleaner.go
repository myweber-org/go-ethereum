
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
}
package main

import (
    "fmt"
    "strings"
)

// TrimSpaces removes leading and trailing whitespace from each string in the slice.
// It returns a new slice with the trimmed strings.
func TrimSpaces(input []string) []string {
    trimmed := make([]string, len(input))
    for i, s := range input {
        trimmed[i] = strings.TrimSpace(s)
    }
    return trimmed
}

func main() {
    data := []string{"  apple ", "banana  ", "  cherry  ", "date"}
    result := TrimSpaces(data)
    fmt.Println("Original:", data)
    fmt.Println("Trimmed:", result)
}