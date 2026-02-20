
package data

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type UserData struct {
	ID       int
	Email    string
	Username string
	Active   bool
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(strings.ToLower(username))
}

func ProcessUserData(user UserData) (UserData, error) {
	if err := ValidateEmail(user.Email); err != nil {
		return UserData{}, err
	}

	normalizedUsername := NormalizeUsername(user.Username)
	if normalizedUsername == "" {
		return UserData{}, errors.New("username cannot be empty after normalization")
	}

	return UserData{
		ID:       user.ID,
		Email:    strings.ToLower(user.Email),
		Username: normalizedUsername,
		Active:   user.Active,
	}, nil
}

func FilterActiveUsers(users []UserData) []UserData {
	var activeUsers []UserData
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func GenerateUserReport(users []UserData) map[string]int {
	report := make(map[string]int)
	for _, user := range users {
		report["total"]++
		if user.Active {
			report["active"]++
		}
	}
	return report
}