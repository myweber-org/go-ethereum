package main

import (
	"fmt"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUser(data UserData) (bool, []string) {
	var errors []string

	if strings.TrimSpace(data.Username) == "" {
		errors = append(errors, "username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		errors = append(errors, "invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		errors = append(errors, "age must be between 0 and 150")
	}

	return len(errors) == 0, errors
}

func TransformUsername(data *UserData) {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
}

func ProcessUserInput(username, email string, age int) (UserData, []string) {
	user := UserData{
		Username: username,
		Email:    email,
		Age:      age,
	}

	TransformUsername(&user)
	valid, errors := ValidateUser(user)

	if !valid {
		return UserData{}, errors
	}

	return user, nil
}

func main() {
	user, errs := ProcessUserInput("  JohnDoe  ", "john@example.com", 25)
	if errs != nil {
		fmt.Println("Validation errors:", errs)
		return
	}
	fmt.Printf("Processed user: %+v\n", user)
}
package main

import (
	"errors"
	"strings"
)

type UserData struct {
	Username string
	Email    string
	Age      int
}

func ValidateUserData(data UserData) error {
	if strings.TrimSpace(data.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if !strings.Contains(data.Email, "@") {
		return errors.New("invalid email format")
	}
	if data.Age < 0 || data.Age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}

func TransformUsername(data UserData) UserData {
	data.Username = strings.ToLower(strings.TrimSpace(data.Username))
	return data
}

func ProcessUserInput(rawUsername string, rawEmail string, rawAge int) (UserData, error) {
	userData := UserData{
		Username: rawUsername,
		Email:    rawEmail,
		Age:      rawAge,
	}

	userData = TransformUsername(userData)

	if err := ValidateUserData(userData); err != nil {
		return UserData{}, err
	}

	return userData, nil
}