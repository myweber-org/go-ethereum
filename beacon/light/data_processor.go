
package main

import "fmt"

func calculateMovingAverage(data []float64, windowSize int) []float64 {
    if len(data) == 0 || windowSize <= 0 || windowSize > len(data) {
        return []float64{}
    }

    result := make([]float64, len(data)-windowSize+1)
    var sum float64

    for i := 0; i < windowSize; i++ {
        sum += data[i]
    }
    result[0] = sum / float64(windowSize)

    for i := windowSize; i < len(data); i++ {
        sum = sum - data[i-windowSize] + data[i]
        result[i-windowSize+1] = sum / float64(windowSize)
    }

    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3

    movingAvg := calculateMovingAverage(sampleData, window)
    fmt.Printf("Moving average (window=%d): %v\n", window, movingAvg)
}