
package main

import (
	"strings"
	"unicode"
)

type DataProcessor struct {
	trimSpaces  bool
	toLowerCase bool
}

func NewDataProcessor(trimSpaces, toLowerCase bool) *DataProcessor {
	return &DataProcessor{
		trimSpaces:  trimSpaces,
		toLowerCase: toLowerCase,
	}
}

func (dp *DataProcessor) Process(input string) string {
	result := input

	if dp.trimSpaces {
		result = strings.TrimSpace(result)
		result = strings.Join(strings.Fields(result), " ")
	}

	if dp.toLowerCase {
		result = strings.ToLower(result)
	}

	return dp.removeSpecialChars(result)
}

func (dp *DataProcessor) removeSpecialChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)
}

func (dp *DataProcessor) BatchProcess(inputs []string) []string {
	results := make([]string, len(inputs))
	for i, input := range inputs {
		results[i] = dp.Process(input)
	}
	return results
}