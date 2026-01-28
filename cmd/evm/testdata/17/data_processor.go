
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ValidateAndParseJSON(rawData []byte) (*UserData, error) {
	var data UserData
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if data.ID <= 0 {
		return nil, fmt.Errorf("invalid ID: must be positive integer")
	}
	if data.Name == "" {
		return nil, fmt.Errorf("name field cannot be empty")
	}
	if data.Email == "" {
		return nil, fmt.Errorf("email field cannot be empty")
	}

	return &data, nil
}

func main() {
	jsonInput := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	parsedData, err := ValidateAndParseJSON([]byte(jsonInput))
	if err != nil {
		log.Fatalf("Error processing data: %v", err)
	}
	fmt.Printf("Successfully parsed user: %+v\n", parsedData)
}