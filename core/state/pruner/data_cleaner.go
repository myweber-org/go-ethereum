package main

import "fmt"

func RemoveDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range input {
		if !seen[item] {
			seen[item] = true
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
}
package main

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
    sample := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
    cleaned := RemoveDuplicates(sample)
    fmt.Println("Original:", sample)
    fmt.Println("Cleaned:", cleaned)
}