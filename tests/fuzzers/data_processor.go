package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Active    bool   `json:"active"`
	Tags      []string `json:"tags"`
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func FilterInactiveUsers(users []UserProfile) []UserProfile {
	var activeUsers []UserProfile
	for _, user := range users {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers
}

func TransformUserData(users []UserProfile) ([]map[string]interface{}, error) {
	var transformed []map[string]interface{}
	
	for _, user := range users {
		if !ValidateEmail(user.Email) {
			return nil, fmt.Errorf("invalid email for user %d: %s", user.ID, user.Email)
		}
		
		transformedUser := map[string]interface{}{
			"user_id":   user.ID,
			"username":  NormalizeUsername(user.Username),
			"email":     strings.ToLower(user.Email),
			"age_group": determineAgeGroup(user.Age),
			"tag_count": len(user.Tags),
			"status":    "active",
		}
		
		if !user.Active {
			transformedUser["status"] = "inactive"
		}
		
		transformed = append(transformed, transformedUser)
	}
	
	return transformed, nil
}

func determineAgeGroup(age int) string {
	switch {
	case age < 18:
		return "minor"
	case age >= 18 && age <= 35:
		return "young_adult"
	case age > 35 && age <= 60:
		return "adult"
	default:
		return "senior"
	}
}

func ProcessUserJSON(jsonData []byte) ([]map[string]interface{}, error) {
	var users []UserProfile
	
	if err := json.Unmarshal(jsonData, &users); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	
	activeUsers := FilterInactiveUsers(users)
	
	transformed, err := TransformUserData(activeUsers)
	if err != nil {
		return nil, err
	}
	
	return transformed, nil
}

func main() {
	jsonData := []byte(`[
		{"id": 1, "username": "JohnDoe ", "email": "john@example.com", "age": 25, "active": true, "tags": ["golang", "backend"]},
		{"id": 2, "username": "JaneSmith", "email": "jane@example.org", "age": 42, "active": false, "tags": ["frontend"]},
		{"id": 3, "username": " BobLee ", "email": "bob@test.com", "age": 17, "active": true, "tags": []},
		{"id": 4, "username": "InvalidUser", "email": "invalid-email", "age": 30, "active": true, "tags": ["test"]}
	]`)
	
	result, err := ProcessUserJSON(jsonData)
	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
		return
	}
	
	for _, user := range result {
		fmt.Printf("%+v\n", user)
	}
}