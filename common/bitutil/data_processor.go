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
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataProcessor struct {
	InputPath  string
	OutputPath string
}

func NewDataProcessor(input, output string) *DataProcessor {
	return &DataProcessor{
		InputPath:  input,
		OutputPath: output,
	}
}

func (dp *DataProcessor) Process() error {
	inputFile, err := os.Open(dp.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dp.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %w", err)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	recordCount := 0
	cleanedCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		recordCount++
		cleanedRecord := dp.cleanRecord(record)
		if dp.isValidRecord(cleanedRecord) {
			if err := writer.Write(cleanedRecord); err != nil {
				return fmt.Errorf("failed to write record: %w", err)
			}
			cleanedCount++
		}
	}

	fmt.Printf("Processed %d records, saved %d valid records\n", recordCount, cleanedCount)
	return nil
}

func (dp *DataProcessor) cleanRecord(record []string) []string {
	cleaned := make([]string, len(record))
	for i, field := range record {
		cleaned[i] = strings.TrimSpace(field)
	}
	return cleaned
}

func (dp *DataProcessor) isValidRecord(record []string) bool {
	for _, field := range record {
		if field == "" {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_processor <input.csv> <output.csv>")
		os.Exit(1)
	}

	processor := NewDataProcessor(os.Args[1], os.Args[2])
	if err := processor.Process(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}