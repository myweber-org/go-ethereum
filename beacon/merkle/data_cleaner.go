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

func (dc *DataCleaner) Process(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", fmt.Errorf("empty input")
	}

	normalized := strings.ToLower(strings.TrimSpace(input))
	if dc.seen[normalized] {
		return "", fmt.Errorf("duplicate entry: %s", input)
	}

	dc.seen[normalized] = true
	return strings.Title(strings.ToLower(input)), nil
}

func (dc *DataCleaner) Reset() {
	dc.seen = make(map[string]bool)
}

func main() {
	cleaner := NewDataCleaner()
	samples := []string{"  apple  ", "APPLE", "banana", "  ", "cherry"}

	for _, sample := range samples {
		result, err := cleaner.Process(sample)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Cleaned: %s -> %s\n", sample, result)
		}
	}
}