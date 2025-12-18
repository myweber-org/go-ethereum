package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values
// using a specified window size. Returns a slice of averages.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return nil
	}

	var result []float64
	var sum float64

	// Calculate the first window's sum
	for i := 0; i < windowSize; i++ {
		sum += data[i]
	}
	result = append(result, sum/float64(windowSize))

	// Slide the window and compute subsequent averages
	for i := windowSize; i < len(data); i++ {
		sum = sum - data[i-windowSize] + data[i]
		result = append(result, sum/float64(windowSize))
	}

	return result
}

func main() {
	// Example usage
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := CalculateMovingAverage(data, window)
	fmt.Printf("Moving averages with window size %d: %v\n", window, averages)
}