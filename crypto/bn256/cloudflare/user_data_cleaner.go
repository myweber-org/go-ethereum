package main

import (
    "regexp"
    "strings"
)

type UserDataCleaner struct{}

func (c *UserDataCleaner) TrimWhitespace(input string) string {
    return strings.TrimSpace(input)
}

func (c *UserDataCleaner) RemoveExtraSpaces(input string) string {
    space := regexp.MustCompile(`\s+`)
    return space.ReplaceAllString(input, " ")
}

func (c *UserDataCleaner) SanitizeEmail(email string) string {
    email = strings.ToLower(c.TrimWhitespace(email))
    emailPattern := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    if emailPattern.MatchString(email) {
        return email
    }
    return ""
}

func (c *UserDataCleaner) ValidateUsername(username string) bool {
    username = c.TrimWhitespace(username)
    if len(username) < 3 || len(username) > 20 {
        return false
    }
    validPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    return validPattern.MatchString(username)
}

func (c *UserDataCleaner) CleanPhoneNumber(phone string) string {
    phone = c.TrimWhitespace(phone)
    digitsOnly := regexp.MustCompile(`\D`)
    cleaned := digitsOnly.ReplaceAllString(phone, "")
    if len(cleaned) == 10 {
        return cleaned
    }
    return ""
}