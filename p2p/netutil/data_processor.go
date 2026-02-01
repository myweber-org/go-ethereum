
package main

import "fmt"

func calculateMovingAverage(data []float64, windowSize int) []float64 {
    if windowSize <= 0 || windowSize > len(data) {
        return nil
    }

    result := make([]float64, 0, len(data)-windowSize+1)
    var sum float64

    for i := 0; i < windowSize; i++ {
        sum += data[i]
    }
    result = append(result, sum/float64(windowSize))

    for i := windowSize; i < len(data); i++ {
        sum = sum - data[i-windowSize] + data[i]
        result = append(result, sum/float64(windowSize))
    }

    return result
}

func main() {
    sampleData := []float64{10.5, 11.2, 10.8, 12.1, 11.9, 12.5, 13.2}
    window := 3

    movingAvg := calculateMovingAverage(sampleData, window)
    fmt.Printf("Moving average (window=%d): %v\n", window, movingAvg)
}