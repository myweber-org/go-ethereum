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
package data_processor

import (
	"regexp"
	"strings"
)

type Processor struct {
	allowedPattern *regexp.Regexp
}

func NewProcessor(allowedPattern string) (*Processor, error) {
	compiled, err := regexp.Compile(allowedPattern)
	if err != nil {
		return nil, err
	}
	return &Processor{allowedPattern: compiled}, nil
}

func (p *Processor) CleanInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return p.allowedPattern.FindString(trimmed)
}

func (p *Processor) ValidateInput(input string) bool {
	return p.allowedPattern.MatchString(input)
}

func (p *Processor) ProcessBatch(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		cleaned := p.CleanInput(input)
		if cleaned != "" {
			results = append(results, cleaned)
		}
	}
	return results
}
package main

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(input, " ")
	
	// Trim spaces from beginning and end
	cleaned = strings.TrimSpace(cleaned)
	
	// Convert to lowercase for consistency
	cleaned = strings.ToLower(cleaned)
	
	return cleaned
}

func NormalizeString(input string) string {
	cleaned := CleanInput(input)
	
	// Remove special characters except alphanumeric and spaces
	re := regexp.MustCompile(`[^a-z0-9\s]`)
	normalized := re.ReplaceAllString(cleaned, "")
	
	return normalized
}

func ProcessData(inputs []string) []string {
	var results []string
	for _, input := range inputs {
		processed := NormalizeString(input)
		if processed != "" {
			results = append(results, processed)
		}
	}
	return results
}package main

import "fmt"

func calculateAverage(numbers []int) float64 {
    if len(numbers) == 0 {
        return 0.0
    }
    
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    
    return float64(sum) / float64(len(numbers))
}

func main() {
    data := []int{10, 20, 30, 40, 50}
    avg := calculateAverage(data)
    fmt.Printf("Average: %.2f\n", avg)
}