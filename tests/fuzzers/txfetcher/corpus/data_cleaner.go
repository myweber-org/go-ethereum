
package main

import (
	"fmt"
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range slice {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func NormalizeString(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func CleanData(data []string) []string {
	cleaned := []string{}
	for _, item := range data {
		normalized := NormalizeString(item)
		cleaned = append(cleaned, normalized)
	}
	return RemoveDuplicates(cleaned)
}

func main() {
	rawData := []string{"  Apple ", "banana", "  apple", "Banana", "Cherry  "}
	cleanedData := CleanData(rawData)
	fmt.Println("Cleaned data:", cleanedData)
}