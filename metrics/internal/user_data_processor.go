
package main

import (
    "errors"
    "regexp"
    "strings"
    "unicode"
)

type UserData struct {
    Username string
    Email    string
    Age      int
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateUserData(data UserData) error {
    if strings.TrimSpace(data.Username) == "" {
        return errors.New("username cannot be empty")
    }
    
    if len(data.Username) < 3 || len(data.Username) > 20 {
        return errors.New("username must be between 3 and 20 characters")
    }
    
    for _, char := range data.Username {
        if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
            return errors.New("username contains invalid characters")
        }
    }
    
    if !emailRegex.MatchString(data.Email) {
        return errors.New("invalid email format")
    }
    
    if data.Age < 13 || data.Age > 120 {
        return errors.New("age must be between 13 and 120")
    }
    
    return nil
}

func SanitizeUsername(username string) string {
    sanitized := strings.TrimSpace(username)
    sanitized = strings.ToLower(sanitized)
    return sanitized
}

func ProcessUserInput(username, email string, age int) (UserData, error) {
    sanitizedUsername := SanitizeUsername(username)
    
    userData := UserData{
        Username: sanitizedUsername,
        Email:    strings.TrimSpace(email),
        Age:      age,
    }
    
    if err := ValidateUserData(userData); err != nil {
        return UserData{}, err
    }
    
    return userData, nil
}