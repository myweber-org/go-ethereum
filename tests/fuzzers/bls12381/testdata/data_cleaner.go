
package main

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
	slice := []int{5, 2, 8, 2, 5, 1, 8, 9}
	uniqueSlice := removeDuplicates(slice)
	fmt.Println("Original:", slice)
	fmt.Println("Unique:", uniqueSlice)
}