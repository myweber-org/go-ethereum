
package main

import "fmt"

func movingAverage(data []float64, windowSize int) []float64 {
    if windowSize <= 0 || len(data) == 0 {
        return []float64{}
    }

    result := make([]float64, 0, len(data)-windowSize+1)
    for i := 0; i <= len(data)-windowSize; i++ {
        sum := 0.0
        for j := 0; j < windowSize; j++ {
            sum += data[i+j]
        }
        result = append(result, sum/float64(windowSize))
    }
    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
    window := 3
    averages := movingAverage(sampleData, window)
    fmt.Printf("Moving averages with window %d: %v\n", window, averages)
}