package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	emailRegex *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return &DataProcessor{emailRegex: regex}
}

func (dp *DataProcessor) SanitizeString(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ProcessUserData(rawEmail, rawName string) (string, string, bool) {
	sanitizedEmail := dp.SanitizeString(rawEmail)
	sanitizedName := dp.SanitizeString(rawName)

	if sanitizedEmail == "" || sanitizedName == "" {
		return sanitizedEmail, sanitizedName, false
	}

	isValid := dp.ValidateEmail(sanitizedEmail)
	return sanitizedEmail, sanitizedName, isValid
}