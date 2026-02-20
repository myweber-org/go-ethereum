
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
)

type UserData struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func ValidateJSON(input string) (*UserData, error) {
    if strings.TrimSpace(input) == "" {
        return nil, errors.New("empty input string")
    }

    var data UserData
    err := json.Unmarshal([]byte(input), &data)
    if err != nil {
        return nil, fmt.Errorf("JSON parsing failed: %w", err)
    }

    if data.Name == "" {
        return nil, errors.New("name field is required")
    }
    if data.Email == "" {
        return nil, errors.New("email field is required")
    }
    if data.Age <= 0 {
        return nil, errors.New("age must be a positive integer")
    }

    return &data, nil
}

func ProcessUserData(jsonStr string) {
    userData, err := ValidateJSON(jsonStr)
    if err != nil {
        fmt.Printf("Validation error: %v\n", err)
        return
    }

    fmt.Printf("Valid user data: Name=%s, Email=%s, Age=%d\n",
        userData.Name, userData.Email, userData.Age)
}

func main() {
    testData := `{"name":"John Doe","email":"john@example.com","age":30}`
    ProcessUserData(testData)

    invalidData := `{"name":"","email":"test@example.com","age":25}`
    ProcessUserData(invalidData)
}