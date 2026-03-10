
package main

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

func NormalizePhone(phone string) string {
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")
	if len(digits) == 10 {
		return digits
	}
	if len(digits) == 11 && digits[0] == '1' {
		return digits[1:]
	}
	return phone
}

func TrimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
package main

import (
	"regexp"
	"strings"
)

type DataProcessor struct {
	allowedPattern *regexp.Regexp
}

func NewDataProcessor() *DataProcessor {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?-]+$`)
	return &DataProcessor{allowedPattern: pattern}
}

func (dp *DataProcessor) SanitizeInput(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	if !dp.allowedPattern.MatchString(trimmed) {
		return "", false
	}

	return trimmed, true
}

func (dp *DataProcessor) ProcessData(data string) (string, error) {
	sanitized, valid := dp.SanitizeInput(data)
	if !valid {
		return "", ErrInvalidInput
	}

	result := strings.ToUpper(sanitized)
	return result, nil
}

var ErrInvalidInput = errors.New("input contains invalid characters")