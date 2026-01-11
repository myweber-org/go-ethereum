
package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

type UserData struct {
    Email     string `json:"email"`
    Username  string `json:"username"`
    Age       int    `json:"age"`
}

func ValidateAndTransform(data []byte) (*UserData, error) {
    var user UserData
    if err := json.Unmarshal(data, &user); err != nil {
        return nil, fmt.Errorf("invalid json format: %w", err)
    }

    user.Email = strings.TrimSpace(strings.ToLower(user.Email))
    user.Username = strings.TrimSpace(user.Username)

    if user.Age < 0 || user.Age > 150 {
        return nil, fmt.Errorf("age %d is out of valid range", user.Age)
    }

    if !strings.Contains(user.Email, "@") {
        return nil, fmt.Errorf("email %s is invalid", user.Email)
    }

    if len(user.Username) == 0 {
        return nil, fmt.Errorf("username cannot be empty")
    }

    return &user, nil
}

func main() {
    jsonData := []byte(`{"email": "TEST@EXAMPLE.COM  ", "username": "  john_doe  ", "age": 25}`)
    processed, err := ValidateAndTransform(jsonData)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Processed: %+v\n", processed)
}