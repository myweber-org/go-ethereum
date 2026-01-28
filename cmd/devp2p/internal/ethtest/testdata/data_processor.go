
package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserProfile struct {
	ID        int
	Email     string
	Username  string
	BirthDate string
	CreatedAt time.Time
}

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
	return strings.ToLower(strings.TrimSpace(username))
}

func CalculateAge(birthDate string) (int, error) {
	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	now := time.Now()
	years := now.Year() - parsedDate.Year()

	if now.YearDay() < parsedDate.YearDay() {
		years--
	}

	if years < 0 {
		return 0, errors.New("birth date cannot be in the future")
	}

	return years, nil
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
	if err := ValidateEmail(profile.Email); err != nil {
		return profile, err
	}

	profile.Username = NormalizeUsername(profile.Username)

	age, err := CalculateAge(profile.BirthDate)
	if err != nil {
		return profile, err
	}

	if age < 13 {
		return profile, errors.New("user must be at least 13 years old")
	}

	return profile, nil
}package data_processor

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func ParseAndValidateJSON(rawData []byte, target interface{}) ([]ValidationError, error) {
	var validationErrors []ValidationError

	if len(rawData) == 0 {
		return append(validationErrors, ValidationError{
			Field:   "payload",
			Message: "empty input data",
		}), nil
	}

	if err := json.Unmarshal(rawData, target); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return append(validationErrors, ValidationError{
				Field:   "syntax",
				Message: fmt.Sprintf("invalid JSON at byte offset %d", syntaxErr.Offset),
			}), nil
		}
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return validationErrors, nil
}

func ProcessUserData(jsonData []byte) (map[string]interface{}, error) {
	var userData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
	}

	validationErrs, err := ParseAndValidateJSON(jsonData, &userData)
	if err != nil {
		return nil, err
	}

	if len(validationErrs) > 0 {
		for _, verr := range validationErrs {
			fmt.Printf("Validation issue: %s\n", verr.Error())
		}
		return nil, fmt.Errorf("data validation failed with %d error(s)", len(validationErrs))
	}

	result := map[string]interface{}{
		"processed":    true,
		"username":     userData.Username,
		"email_domain": extractDomain(userData.Email),
		"age_group":    categorizeAge(userData.Age),
	}

	return result, nil
}

func extractDomain(email string) string {
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			return email[i+1:]
		}
	}
	return ""
}

func categorizeAge(age int) string {
	switch {
	case age < 13:
		return "child"
	case age >= 13 && age < 20:
		return "teen"
	case age >= 20 && age < 65:
		return "adult"
	default:
		return "senior"
	}
}