package main

import (
    "errors"
    "strings"
    "unicode"
)

type UserProfile struct {
    Username string
    Email    string
    Age      int
}

func ValidateProfile(p UserProfile) error {
    if strings.TrimSpace(p.Username) == "" {
        return errors.New("username cannot be empty")
    }
    if len(p.Username) < 3 || len(p.Username) > 20 {
        return errors.New("username must be between 3 and 20 characters")
    }
    for _, r := range p.Username {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
            return errors.New("username can only contain letters, digits, and underscores")
        }
    }

    if !strings.Contains(p.Email, "@") {
        return errors.New("invalid email format")
    }

    if p.Age < 0 || p.Age > 150 {
        return errors.New("age must be between 0 and 150")
    }

    return nil
}

func NormalizeProfile(p UserProfile) UserProfile {
    normalized := p
    normalized.Username = strings.ToLower(strings.TrimSpace(p.Username))
    normalized.Email = strings.ToLower(strings.TrimSpace(p.Email))
    return normalized
}

func ProcessUserProfile(p UserProfile) (UserProfile, error) {
    if err := ValidateProfile(p); err != nil {
        return UserProfile{}, err
    }
    return NormalizeProfile(p), nil
}