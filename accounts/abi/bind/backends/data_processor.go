package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	pattern := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("invalid email format")
	}
	return email, nil
}

func validateUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return "", fmt.Errorf("username must be between 3 and 20 characters")
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, username)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return username, nil
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	email, err := normalizeEmail(data.Email)
	if err != nil {
		return nil, fmt.Errorf("email validation failed: %v", err)
	}
	data.Email = email

	username, err := validateUsername(data.Username)
	if err != nil {
		return nil, fmt.Errorf("username validation failed: %v", err)
	}
	data.Username = username

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150")
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email": "  TEST@Example.COM ", "username": "user_123", "age": 25}`
	processed, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processed)
}
package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DataRecord struct {
	ID        string
	Value     float64
	Timestamp time.Time
	Tags      []string
}

func ValidateRecord(record DataRecord) error {
	if record.ID == "" {
		return errors.New("ID cannot be empty")
	}
	if record.Value < 0 {
		return errors.New("value must be non-negative")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp must be set")
	}
	return nil
}

func TransformRecord(record DataRecord) DataRecord {
	transformed := record
	transformed.Value = record.Value * 1.1
	transformed.Tags = append(record.Tags, "processed")
	return transformed
}

func ProcessRecords(records []DataRecord) ([]DataRecord, error) {
	var processed []DataRecord
	for _, record := range records {
		if err := ValidateRecord(record); err != nil {
			return nil, fmt.Errorf("validation failed for record %s: %w", record.ID, err)
		}
		processed = append(processed, TransformRecord(record))
	}
	return processed, nil
}

func FormatTags(tags []string) string {
	return strings.Join(tags, ", ")
}

func main() {
	records := []DataRecord{
		{
			ID:        "rec001",
			Value:     100.0,
			Timestamp: time.Now(),
			Tags:      []string{"input", "raw"},
		},
		{
			ID:        "rec002",
			Value:     200.0,
			Timestamp: time.Now().Add(-time.Hour),
			Tags:      []string{"input"},
		},
	}

	processed, err := ProcessRecords(records)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}

	for _, rec := range processed {
		fmt.Printf("Processed: ID=%s, Value=%.2f, Tags=[%s]\n",
			rec.ID, rec.Value, FormatTags(rec.Tags))
	}
}