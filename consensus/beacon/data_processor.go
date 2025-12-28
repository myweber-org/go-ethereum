
package data_processor

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	stripSpaces   bool
	removeSpecial bool
}

func NewDataProcessor(stripSpaces, removeSpecial bool) *DataProcessor {
	return &DataProcessor{
		stripSpaces:   stripSpaces,
		removeSpecial: removeSpecial,
	}
}

func (dp *DataProcessor) Process(input string) string {
	result := input

	if dp.stripSpaces {
		result = strings.TrimSpace(result)
	}

	if dp.removeSpecial {
		re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
		result = re.ReplaceAllString(result, "")
	}

	return result
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (dp *DataProcessor) CountWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}