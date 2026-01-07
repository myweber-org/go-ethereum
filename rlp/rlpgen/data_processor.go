
package data_processor

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	whitespaceRegex *regexp.Regexp
	emailRegex      *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		whitespaceRegex: regexp.MustCompile(`\s+`),
		emailRegex:      regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (dp *DataProcessor) CleanString(input string) string {
	trimmed := strings.TrimSpace(input)
	return dp.whitespaceRegex.ReplaceAllString(trimmed, " ")
}

func (dp *DataProcessor) ValidateEmail(email string) bool {
	return dp.emailRegex.MatchString(email)
}

func (dp *DataProcessor) ExtractDomain(email string) (string, bool) {
	if !dp.ValidateEmail(email) {
		return "", false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}package main

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	trimmed := strings.TrimSpace(username)
	return strings.ToLower(trimmed)
}

func SanitizeInput(input string) string {
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}

func TransformPhoneNumber(phone string) (string, error) {
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")

	if len(digits) < 10 {
		return "", errors.New("phone number too short")
	}

	if len(digits) == 10 {
		return "+1" + digits, nil
	}

	if len(digits) == 11 && strings.HasPrefix(digits, "1") {
		return "+" + digits, nil
	}

	if len(digits) > 11 {
		return "", errors.New("phone number too long")
	}

	return "", errors.New("invalid phone number format")
}