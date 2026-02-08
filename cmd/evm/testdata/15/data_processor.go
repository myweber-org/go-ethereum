package main

import (
	"regexp"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Comments string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func sanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`[<>"'&]`)
	return re.ReplaceAllString(input, "")
}

func validateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func processUserData(data UserData) (UserData, error) {
	data.Username = sanitizeInput(data.Username)
	data.Email = sanitizeInput(data.Email)
	data.Comments = sanitizeInput(data.Comments)

	if !validateEmail(data.Email) {
		return data, &InvalidEmailError{Email: data.Email}
	}

	if len(data.Username) < 3 || len(data.Username) > 50 {
		return data, &InvalidUsernameError{Username: data.Username}
	}

	return data, nil
}

type InvalidEmailError struct {
	Email string
}

func (e *InvalidEmailError) Error() string {
	return "Invalid email format: " + e.Email
}

type InvalidUsernameError struct {
	Username string
}

func (e *InvalidUsernameError) Error() string {
	return "Username must be between 3 and 50 characters: " + e.Username
}
package data

import (
	"errors"
	"strings"
	"time"
)

type Record struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

var (
	ErrInvalidID    = errors.New("invalid record ID")
	ErrInvalidValue = errors.New("value out of valid range")
	ErrEmptyTags    = errors.New("record must have at least one tag")
)

func ValidateRecord(r Record) error {
	if len(r.ID) == 0 || len(r.ID) > 100 {
		return ErrInvalidID
	}

	if r.Value < 0 || r.Value > 10000 {
		return ErrInvalidValue
	}

	if len(r.Tags) == 0 {
		return ErrEmptyTags
	}

	return nil
}

func NormalizeTags(tags []string) []string {
	uniqueTags := make(map[string]bool)
	var result []string

	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !uniqueTags[normalized] {
			uniqueTags[normalized] = true
			result = append(result, normalized)
		}
	}

	return result
}

func CalculateAverage(records []Record) (float64, error) {
	if len(records) == 0 {
		return 0, errors.New("cannot calculate average of empty slice")
	}

	var sum float64
	validCount := 0

	for _, r := range records {
		if err := ValidateRecord(r); err == nil {
			sum += r.Value
			validCount++
		}
	}

	if validCount == 0 {
		return 0, errors.New("no valid records to calculate average")
	}

	return sum / float64(validCount), nil
}

func FilterByTag(records []Record, tag string) []Record {
	var filtered []Record
	targetTag := strings.ToLower(strings.TrimSpace(tag))

	for _, r := range records {
		for _, t := range r.Tags {
			if strings.ToLower(t) == targetTag {
				filtered = append(filtered, r)
				break
			}
		}
	}

	return filtered
}

func TransformRecords(records []Record, multiplier float64) []Record {
	transformed := make([]Record, len(records))

	for i, r := range records {
		transformed[i] = Record{
			ID:        r.ID,
			Value:     r.Value * multiplier,
			Timestamp: r.Timestamp,
			Tags:      NormalizeTags(r.Tags),
		}
	}

	return transformed
}