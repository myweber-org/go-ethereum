
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
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if data.ID <= 0 {
		return nil, fmt.Errorf("invalid ID: must be positive integer")
	}
	if data.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if data.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	return &data, nil
}

func main() {
	jsonStr := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
	parsedData, err := ValidateAndParseJSON([]byte(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed data: %+v\n", parsedData)
}