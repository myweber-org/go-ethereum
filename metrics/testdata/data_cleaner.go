
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
}package main

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
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	normalizeCase bool
	trimSpaces    bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		normalizeCase: true,
		trimSpaces:    true,
	}
}

func (dc *DataCleaner) NormalizeString(input string) string {
	result := input

	if dc.trimSpaces {
		result = strings.TrimSpace(result)
	}

	if dc.normalizeCase {
		result = strings.ToLower(result)
	}

	return result
}

func (dc *DataCleaner) DeduplicateStrings(strings []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, s := range strings {
		normalized := dc.NormalizeString(s)
		if _, exists := seen[normalized]; !exists {
			seen[normalized] = struct{}{}
			result = append(result, normalized)
		}
	}

	return result
}

func main() {
	cleaner := NewDataCleaner()

	data := []string{
		"  Apple  ",
		"apple",
		" BANANA ",
		"banana",
		"Cherry",
		"cherry ",
	}

	fmt.Println("Original data:", data)
	cleaned := cleaner.DeduplicateStrings(data)
	fmt.Println("Cleaned data:", cleaned)
}