package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func cleanCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		cleanedRecord := make([]string, len(record))
		for i, field := range record {
			cleanedRecord[i] = strings.TrimSpace(field)
		}

		if err := writer.Write(cleanedRecord); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := cleanCSV(inputFile, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully cleaned CSV: %s -> %s\n", inputFile, outputFile)
}package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	normalizeCase bool
	trimSpaces    bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		normalizeCase: true,
		trimSpaces:    true,
	}
}

func (dc *DataCleaner) NormalizeString(input string) string {
	result := input
	if dc.trimSpaces {
		result = strings.TrimSpace(result)
	}
	if dc.normalizeCase {
		result = strings.ToLower(result)
	}
	return result
}

func (dc *DataCleaner) DeduplicateStrings(items []string) []string {
	seen := make(map[string]bool)
	var unique []string
	for _, item := range items {
		normalized := dc.NormalizeString(item)
		if !seen[normalized] {
			seen[normalized] = true
			unique = append(unique, normalized)
		}
	}
	return unique
}

func main() {
	cleaner := NewDataCleaner()
	data := []string{"  Apple", "apple ", "BANANA", "banana", "  Cherry  "}
	cleaned := cleaner.DeduplicateStrings(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
}
package main

import (
	"fmt"
	"strings"
)

type DataCleaner struct {
	seen map[string]bool
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		seen: make(map[string]bool),
	}
}

func (dc *DataCleaner) RemoveDuplicates(items []string) []string {
	var unique []string
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if !dc.seen[normalized] && dc.isValid(item) {
			dc.seen[normalized] = true
			unique = append(unique, item)
		}
	}
	return unique
}

func (dc *DataCleaner) isValid(item string) bool {
	return len(item) > 0 && len(item) < 100
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	
	data := []string{
		"apple",
		"Apple",
		"banana",
		"  banana  ",
		"cherry",
		"",
		"cherry",
	}
	
	cleaned := cleaner.RemoveDuplicates(data)
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cleaned: %v\n", cleaned)
	
	cleaner.Reset()
	
	moreData := []string{"grape", "Grape", "grape"}
	moreCleaned := cleaner.RemoveDuplicates(moreData)
	fmt.Printf("More cleaned: %v\n", moreCleaned)
}