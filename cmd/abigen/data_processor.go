
package main

import (
	"regexp"
	"strings"
)

// DataProcessor handles cleaning and normalization of string data
type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
}

// NewDataProcessor creates a new DataProcessor instance
func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
	}
}

// CleanString removes extra whitespace and trims the input
func (dp *DataProcessor) CleanString(input string) string {
	if input == "" {
		return input
	}
	
	trimmed := strings.TrimSpace(input)
	cleaned := dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
	return cleaned
}

// NormalizeString converts string to lowercase and cleans it
func (dp *DataProcessor) NormalizeString(input string) string {
	cleaned := dp.CleanString(input)
	return strings.ToLower(cleaned)
}

// ProcessBatch processes multiple strings and returns cleaned results
func (dp *DataProcessor) ProcessBatch(inputs []string) []string {
	results := make([]string, len(inputs))
	for i, input := range inputs {
		results[i] = dp.NormalizeString(input)
	}
	return results
}