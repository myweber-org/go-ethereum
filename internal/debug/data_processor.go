
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// ValidateJSONString checks if a string is valid JSON.
func ValidateJSONString(s string) bool {
    var js interface{}
    return json.Unmarshal([]byte(s), &js) == nil
}

// PrettyPrintJSON takes a JSON string and prints it with indentation.
func PrettyPrintJSON(jsonStr string) error {
    if !ValidateJSONString(jsonStr) {
        return fmt.Errorf("invalid JSON string")
    }

    var data interface{}
    if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
        return err
    }

    prettyJSON, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }

    fmt.Println(string(prettyJSON))
    return nil
}

// ExtractJSONField attempts to extract a top-level field from a JSON string.
func ExtractJSONField(jsonStr, field string) (string, error) {
    if !ValidateJSONString(jsonStr) {
        return "", fmt.Errorf("invalid JSON string")
    }

    var result map[string]interface{}
    if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
        return "", err
    }

    value, exists := result[field]
    if !exists {
        return "", fmt.Errorf("field '%s' not found", field)
    }

    // Convert the value back to a JSON string
    valueBytes, err := json.Marshal(value)
    if err != nil {
        return "", err
    }
    return string(valueBytes), nil
}

// IsJSONObject checks if the given string is a JSON object (starts with '{').
func IsJSONObject(s string) bool {
    trimmed := strings.TrimSpace(s)
    return strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")
}

func main() {
    // Example usage
    testJSON := `{"name":"Alice","age":30,"active":true}`
    fmt.Println("Valid JSON?", ValidateJSONString(testJSON))

    fmt.Println("\nPretty Print:")
    PrettyPrintJSON(testJSON)

    field, err := ExtractJSONField(testJSON, "name")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("\nExtracted 'name' field:", field)
    }

    fmt.Println("\nIs JSON Object?", IsJSONObject(testJSON))
}