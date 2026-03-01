package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	data := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	cleaned := RemoveDuplicates(data)
	fmt.Println("Original:", data)
	fmt.Println("Cleaned:", cleaned)
}package datautils

import "sort"

func RemoveDuplicates[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(input))

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesSorted[T comparable](input []T) []T {
	if len(input) == 0 {
		return input
	}

	sort.Slice(input, func(i, j int) bool {
		// Convert to string for comparison to satisfy sort.Interface
		// This is a limitation for generic comparable types
		return false // Just a placeholder, actual implementation would need type assertion
	})

	result := input[:1]
	for i := 1; i < len(input); i++ {
		if input[i] != input[i-1] {
			result = append(result, input[i])
		}
	}

	return result
}