
package main

import (
    "encoding/json"
    "fmt"
    "regexp"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func ValidateUser(u User) error {
    if u.ID <= 0 {
        return fmt.Errorf("invalid user ID: %d", u.ID)
    }
    if len(u.Username) < 3 || len(u.Username) > 20 {
        return fmt.Errorf("username must be between 3 and 20 characters")
    }
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(u.Email) {
        return fmt.Errorf("invalid email format: %s", u.Email)
    }
    return nil
}

func ProcessUserData(jsonData []byte) (*User, error) {
    var user User
    err := json.Unmarshal(jsonData, &user)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    err = ValidateUser(user)
    if err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    return &user, nil
}

func main() {
    data := []byte(`{"id": 123, "username": "john_doe", "email": "john@example.com"}`)
    user, err := ProcessUserData(data)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Processed user: %+v\n", user)
}