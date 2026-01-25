
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
package main

import (
    "regexp"
    "strings"
)

type DataProcessor struct {
    stripPattern *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        stripPattern: regexp.MustCompile(`[^a-zA-Z0-9\s\-_]`),
    }
}

func (dp *DataProcessor) CleanInput(input string) string {
    cleaned := dp.stripPattern.ReplaceAllString(input, "")
    return strings.TrimSpace(cleaned)
}

func (dp *DataProcessor) ValidateLength(input string, min, max int) bool {
    length := len(input)
    return length >= min && length <= max
}

func (dp *DataProcessor) Process(input string, minLen, maxLen int) (string, bool) {
    cleaned := dp.CleanInput(input)
    isValid := dp.ValidateLength(cleaned, minLen, maxLen)
    return cleaned, isValid
}