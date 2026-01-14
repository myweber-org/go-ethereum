
package main

import (
    "errors"
    "strings"
    "time"
)

type User struct {
    ID        int
    Username  string
    Email     string
    CreatedAt time.Time
    Active    bool
}

func ValidateUser(user User) error {
    if user.Username == "" {
        return errors.New("username cannot be empty")
    }
    if len(user.Username) < 3 {
        return errors.New("username must be at least 3 characters")
    }
    if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email format")
    }
    if user.ID <= 0 {
        return errors.New("invalid user ID")
    }
    return nil
}

func NormalizeUsername(username string) string {
    return strings.ToLower(strings.TrimSpace(username))
}

func FormatUserDisplay(user User) string {
    status := "inactive"
    if user.Active {
        status = "active"
    }
    return strings.Join([]string{
        "User: " + user.Username,
        "Email: " + user.Email,
        "Status: " + status,
        "Created: " + user.CreatedAt.Format("2006-01-02"),
    }, " | ")
}

func FilterActiveUsers(users []User) []User {
    var activeUsers []User
    for _, user := range users {
        if user.Active {
            activeUsers = append(activeUsers, user)
        }
    }
    return activeUsers
}

func CalculateAverageAge(users []User, currentTime time.Time) float64 {
    if len(users) == 0 {
        return 0
    }
    
    var totalDays float64
    for _, user := range users {
        age := currentTime.Sub(user.CreatedAt).Hours() / 24
        totalDays += age
    }
    
    return totalDays / float64(len(users))
}