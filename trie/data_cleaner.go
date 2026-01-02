package main

import (
	"fmt"
	"strings"
)

func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func CleanData(data []string) []string {
	normalized := make([]string, len(data))
	for i, v := range data {
		normalized[i] = NormalizeString(v)
	}
	return RemoveDuplicates(normalized)
}

func main() {
	sampleData := []string{" Apple ", "banana", "apple", " Banana ", "Cherry"}
	cleaned := CleanData(sampleData)
	fmt.Println("Cleaned data:", cleaned)
}