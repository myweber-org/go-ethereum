
package main

import (
	"fmt"
)

// CalculateMovingAverage computes the moving average of a slice of float64 values.
// It takes a slice of data and a window size as input.
// Returns a slice containing the moving averages or an empty slice if window size is invalid.
func CalculateMovingAverage(data []float64, windowSize int) []float64 {
	if windowSize <= 0 || windowSize > len(data) {
		return []float64{}
	}

	result := make([]float64, len(data)-windowSize+1)
	for i := 0; i <= len(data)-windowSize; i++ {
		sum := 0.0
		for j := i; j < i+windowSize; j++ {
			sum += data[j]
		}
		result[i] = sum / float64(windowSize)
	}
	return result
}

func main() {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	window := 3
	averages := CalculateMovingAverage(data, window)
	fmt.Printf("Moving averages with window size %d: %v\n", window, averages)
}