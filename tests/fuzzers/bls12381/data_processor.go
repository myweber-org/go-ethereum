package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitized := re.ReplaceAllString(input, "")
	return strings.TrimSpace(sanitized)
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}package main

import "fmt"

func calculateMovingAverage(data []float64, windowSize int) []float64 {
    if len(data) == 0 || windowSize <= 0 || windowSize > len(data) {
        return []float64{}
    }

    result := make([]float64, len(data)-windowSize+1)
    for i := 0; i <= len(data)-windowSize; i++ {
        sum := 0.0
        for j := 0; j < windowSize; j++ {
            sum += data[i+j]
        }
        result[i] = sum / float64(windowSize)
    }
    return result
}

func main() {
    sampleData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
    window := 3
    averages := calculateMovingAverage(sampleData, window)
    
    fmt.Printf("Original data: %v\n", sampleData)
    fmt.Printf("Moving averages (window=%d): %v\n", window, averages)
}