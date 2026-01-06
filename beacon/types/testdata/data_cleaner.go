
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
	slice := []string{"apple", "banana", "apple", "orange", "banana", "grape"}
	unique := RemoveDuplicates(slice)
	fmt.Println("Original:", slice)
	fmt.Println("Unique:", unique)
}