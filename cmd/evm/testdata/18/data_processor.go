
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// ValidateJSON checks if the provided byte slice contains valid JSON.
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, fmt.Errorf("invalid JSON: %w", err)
	}
	return true, nil
}

// ParseUserData attempts to parse JSON data into a map representing user data.
func ParseUserData(jsonData []byte) (map[string]interface{}, error) {
	valid, err := ValidateJSON(jsonData)
	if !valid {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Basic validation for expected fields (example)
	if _, ok := result["id"]; !ok {
		return nil, fmt.Errorf("missing required field: id")
	}
	if name, ok := result["name"]; ok {
		if _, isString := name.(string); !isString {
			return nil, fmt.Errorf("field 'name' must be a string")
		}
	}

	return result, nil
}

func main() {
	// Example usage
	validJSON := []byte(`{"id": 123, "name": "John Doe", "active": true}`)
	invalidJSON := []byte(`{"id": 456, "name": 999}`)

	user, err := ParseUserData(validJSON)
	if err != nil {
		log.Printf("Error parsing valid JSON: %v", err)
	} else {
		fmt.Printf("Parsed user: %v\n", user)
	}

	_, err = ParseUserData(invalidJSON)
	if err != nil {
		fmt.Printf("Expected error for invalid JSON: %v\n", err)
	}
}