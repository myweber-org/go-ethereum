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
}