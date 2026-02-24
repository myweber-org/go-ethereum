package datautils

import "sort"

func RemoveDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicatesSorted[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	sort.Slice(slice, func(i, j int) bool {
		// Simple comparison for sorting
		return false // Default case, actual implementation would compare elements
	})

	result := slice[:1]
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[i-1] {
			result = append(result, slice[i])
		}
	}
	return result
}package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DataCleaner struct {
	InputPath  string
	OutputPath string
	Delimiter  rune
}

func NewDataCleaner(input, output string) *DataCleaner {
	return &DataCleaner{
		InputPath:  input,
		OutputPath: output,
		Delimiter:  ',',
	}
}

func (dc *DataCleaner) Clean() error {
	inputFile, err := os.Open(dc.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(dc.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	reader := csv.NewReader(inputFile)
	reader.Comma = dc.Delimiter
	writer := csv.NewWriter(outputFile)
	writer.Comma = dc.Delimiter
	defer writer.Flush()

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	seen := make(map[string]bool)
	recordCount := 0
	duplicateCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		recordKey := strings.Join(record, "|")
		if seen[recordKey] {
			duplicateCount++
			continue
		}

		seen[recordKey] = true
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
		recordCount++
	}

	fmt.Printf("Processed %d records, removed %d duplicates\n", recordCount, duplicateCount)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: data_cleaner <input.csv> <output.csv>")
		os.Exit(1)
	}

	cleaner := NewDataCleaner(os.Args[1], os.Args[2])
	if err := cleaner.Clean(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}