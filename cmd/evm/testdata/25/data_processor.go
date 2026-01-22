package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONValidator holds configuration for validation
type JSONValidator struct {
	RequiredFields []string
}

// ParseAndValidate attempts to parse a JSON string and validate required fields.
// Returns parsed data as map[string]interface{} and an error if validation fails.
func (v *JSONValidator) ParseAndValidate(rawData string) (map[string]interface{}, error) {
	var data map[string]interface{}

	decoder := json.NewDecoder(strings.NewReader(rawData))
	decoder.UseNumber()

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, field := range v.RequiredFields {
		if _, exists := data[field]; !exists {
			return nil, fmt.Errorf("missing required field: %s", field)
		}
	}

	return data, nil
}

// NewValidator creates a new JSONValidator with the given required fields.
func NewValidator(requiredFields []string) *JSONValidator {
	return &JSONValidator{
		RequiredFields: requiredFields,
	}
}

func main() {
	validator := NewValidator([]string{"id", "name"})

	sampleJSON := `{"id": 123, "name": "Test Item", "active": true}`

	result, err := validator.ParseAndValidate(sampleJSON)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	fmt.Printf("Validated data: %v\n", result)
}package main

import (
	"errors"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string
	Username string
	Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(data.Email) {
		return errors.New("invalid email format")
	}
	if len(data.Username) < 3 || len(data.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if data.Age < 18 || data.Age > 120 {
		return errors.New("age must be between 18 and 120")
	}
	return nil
}

func TransformUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func ProcessUserInput(email, username string, age int) (UserData, error) {
	transformedUsername := TransformUsername(username)
	userData := UserData{
		Email:    strings.TrimSpace(email),
		Username: transformedUsername,
		Age:      age,
	}
	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}
	return userData, nil
}