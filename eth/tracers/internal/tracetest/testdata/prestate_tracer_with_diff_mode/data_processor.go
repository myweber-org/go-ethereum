
package main

import (
    "fmt"
    "strings"
    "unicode"
)

type UserProfile struct {
    Username string
    Email    string
    Age      int
}

func NormalizeUsername(username string) string {
    return strings.ToLower(strings.TrimSpace(username))
}

func ValidateEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidateAge(age int) bool {
    return age >= 0 && age <= 120
}

func ProcessUserProfile(profile UserProfile) (UserProfile, error) {
    normalizedUsername := NormalizeUsername(profile.Username)
    if normalizedUsername == "" {
        return UserProfile{}, fmt.Errorf("username cannot be empty")
    }

    if !ValidateEmail(profile.Email) {
        return UserProfile{}, fmt.Errorf("invalid email format")
    }

    if !ValidateAge(profile.Age) {
        return UserProfile{}, fmt.Errorf("age must be between 0 and 120")
    }

    return UserProfile{
        Username: normalizedUsername,
        Email:    strings.ToLower(strings.TrimSpace(profile.Email)),
        Age:      profile.Age,
    }, nil
}

func IsAlphanumeric(str string) bool {
    for _, char := range str {
        if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
            return false
        }
    }
    return true
}

func main() {
    testProfile := UserProfile{
        Username: "  TestUser123  ",
        Email:    "test@example.com",
        Age:      25,
    }

    processed, err := ProcessUserProfile(testProfile)
    if err != nil {
        fmt.Printf("Error processing profile: %v\n", err)
        return
    }

    fmt.Printf("Processed profile: %+v\n", processed)
    fmt.Printf("Username is alphanumeric: %v\n", IsAlphanumeric(processed.Username))
}