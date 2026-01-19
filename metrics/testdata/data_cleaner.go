
package main

import (
	"fmt"
	"strings"
)

func CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	return lower
}

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func main() {
	data := []string{" Apple", "banana ", "Apple", "  BANANA", "Cherry"}
	cleaned := []string{}

	for _, item := range data {
		cleaned = append(cleaned, CleanString(item))
	}

	unique := RemoveDuplicates(cleaned)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
	fmt.Println("Unique:", unique)
}