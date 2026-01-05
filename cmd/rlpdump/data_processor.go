package main

import (
	"fmt"
	"strings"
	"unicode"
)

type UserData struct {
	Username string
	Email    string
}

func NormalizeUsername(username string) string {
	trimmed := strings.TrimSpace(username)
	var result strings.Builder
	for _, r := range trimmed {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return strings.Contains(parts[1], ".")
}

func ProcessUserInput(username, email string) (*UserData, error) {
	normalizedUsername := NormalizeUsername(username)
	if len(normalizedUsername) < 3 {
		return nil, fmt.Errorf("username must be at least 3 characters after normalization")
	}

	if !ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	return &UserData{
		Username: normalizedUsername,
		Email:    strings.ToLower(strings.TrimSpace(email)),
	}, nil
}

func main() {
	user, err := ProcessUserInput("  John_Doe-123  ", "john@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}