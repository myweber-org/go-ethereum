
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// DataPayload represents a generic structure for incoming JSON data
type DataPayload struct {
    ID    string          `json:"id"`
    Value json.RawMessage `json:"value"`
    Tags  []string        `json:"tags"`
}

// ValidateJSON checks if the provided string is valid JSON and unmarshals it into DataPayload
func ValidateJSON(input string) (*DataPayload, error) {
    var payload DataPayload
    decoder := json.NewDecoder(strings.NewReader(input))
    decoder.DisallowUnknownFields()

    if err := decoder.Decode(&payload); err != nil {
        return nil, fmt.Errorf("invalid JSON structure: %w", err)
    }

    if payload.ID == "" {
        return nil, fmt.Errorf("field 'id' is required and cannot be empty")
    }

    if len(payload.Tags) > 10 {
        return nil, fmt.Errorf("field 'tags' cannot contain more than 10 items")
    }

    return &payload, nil
}

func main() {
    testJSON := `{"id": "user_123", "value": {"score": 95}, "tags": ["go", "backend"]}`
    result, err := ValidateJSON(testJSON)
    if err != nil {
        fmt.Printf("Validation failed: %v\n", err)
        return
    }
    fmt.Printf("Valid payload: ID=%s, Tags=%v\n", result.ID, result.Tags)
}