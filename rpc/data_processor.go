
package data_processor

import (
	"regexp"
	"strings"
	"unicode"
)

func CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	normalized := normalizeSpaces(trimmed)
	return removeSpecialChars(normalized)
}

func normalizeSpaces(s string) string {
	spaceRegex := regexp.MustCompile(`\s+`)
	return spaceRegex.ReplaceAllString(s, " ")
}

func removeSpecialChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, s)
}

func Tokenize(s string) []string {
	cleaned := CleanInput(s)
	if cleaned == "" {
		return []string{}
	}
	return strings.Split(cleaned, " ")
}
package main

import (
	"errors"
	"strings"
)

func ProcessUserInput(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("input cannot be empty")
	}

	processed := strings.ToLower(strings.TrimSpace(input))
	processed = strings.ReplaceAll(processed, "  ", " ")

	return processed, nil
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
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if !validateEmail(data.Email) {
		return nil, fmt.Errorf("invalid email format: %s", data.Email)
	}

	data.Username = sanitizeUsername(data.Username)

	if data.Age < 0 || data.Age > 150 {
		return nil, fmt.Errorf("age out of valid range: %d", data.Age)
	}

	return &data, nil
}

func main() {
	rawJSON := `{"email":"test@example.com","username":"  john_doe  ","age":25}`
	processedData, err := processUserData([]byte(rawJSON))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed Data: %+v\n", processedData)
}