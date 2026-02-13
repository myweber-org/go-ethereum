
package main

import "fmt"

func RemoveDuplicates(nums []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, num := range nums {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    
    return result
}

func main() {
    sampleData := []int{1, 2, 2, 3, 4, 4, 5, 1, 6}
    cleaned := RemoveDuplicates(sampleData)
    fmt.Printf("Original: %v\n", sampleData)
    fmt.Printf("Cleaned: %v\n", cleaned)
}