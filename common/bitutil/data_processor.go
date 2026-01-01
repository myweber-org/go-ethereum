package main

import (
	"fmt"
	"math"
)

// CalculateMovingAverage computes the simple moving average for a slice of float64 values.
// windowSize specifies the number of data points to include in each average calculation.
// Returns a slice of moving averages and an error if the window size is invalid.
func CalculateMovingAverage(data []float64, windowSize int) ([]float64, error) {
	if windowSize <= 0 {
		return nil, fmt.Errorf("window size must be positive, got %d", windowSize)
	}
	if len(data) < windowSize {
		return nil, fmt.Errorf("data length %d is less than window size %d", len(data), windowSize)
	}

	var result []float64
	for i := 0; i <= len(data)-windowSize; i++ {
		sum := 0.0
		for j := i; j < i+windowSize; j++ {
			sum += data[j]
		}
		average := sum / float64(windowSize)
		// Round to two decimal places for cleaner output
		rounded := math.Round(average*100) / 100
		result = append(result, rounded)
	}
	return result, nil
}

func main() {
	// Example usage
	stockPrices := []float64{45.12, 46.25, 47.80, 46.95, 48.10, 49.35, 50.20, 49.80}
	window := 3

	averages, err := CalculateMovingAverage(stockPrices, window)
	if err != nil {
		fmt.Printf("Error calculating moving average: %v\n", err)
		return
	}

	fmt.Printf("Original data: %v\n", stockPrices)
	fmt.Printf("%d-day moving averages: %v\n", window, averages)
}