
package data_processor

import (
	"regexp"
	"strings"
)

type DataCleaner struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dc *DataCleaner) NormalizeWhitespace(input string) string {
	trimmed := strings.TrimSpace(input)
	return dc.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dc *DataCleaner) ValidateEmail(email string) bool {
	return dc.emailRegex.MatchString(email)
}

func (dc *DataCleaner) RemoveSpecialChars(input string, keep string) string {
	pattern := "[^a-zA-Z0-9" + regexp.QuoteMeta(keep) + "]+"
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, "")
}

func (dc *DataCleaner) TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	return input[:maxLength-3] + "..."
}