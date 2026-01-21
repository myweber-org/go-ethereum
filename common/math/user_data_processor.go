package main

import (
	"regexp"
	"strings"
)

func SanitizeUsername(input string) (string, error) {
	if input == "" {
		return "", ErrEmptyInput
	}

	trimmed := strings.TrimSpace(input)

	pattern := `^[a-zA-Z0-9_\-\.]+$`
	matched, err := regexp.MatchString(pattern, trimmed)
	if err != nil {
		return "", err
	}

	if !matched {
		return "", ErrInvalidCharacters
	}

	if len(trimmed) < 3 || len(trimmed) > 32 {
		return "", ErrLengthViolation
	}

	return trimmed, nil
}

var (
	ErrEmptyInput        = errors.New("input string cannot be empty")
	ErrInvalidCharacters = errors.New("username contains invalid characters")
	ErrLengthViolation   = errors.New("username must be between 3 and 32 characters")
)