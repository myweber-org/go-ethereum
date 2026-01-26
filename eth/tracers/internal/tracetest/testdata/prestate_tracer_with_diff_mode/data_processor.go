
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
}package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func sanitizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func processUserData(rawData []byte) (*UserData, error) {
	var data UserData
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	if !validateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	data.Username = sanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range")
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	processed, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed data: %+v\n", processed)
}