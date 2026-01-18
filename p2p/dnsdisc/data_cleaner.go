package datautils

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
}package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range input {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func main() {
	slice := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	uniqueSlice := RemoveDuplicates(slice)
	fmt.Println("Original:", slice)
	fmt.Println("Unique:", uniqueSlice)
}