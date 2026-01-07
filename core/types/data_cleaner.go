
package main

import (
    "fmt"
    "strings"
)

func DeduplicateStrings(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
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
    if !strings.Contains(parts[1], ".") {
        return false
    }
    return true
}

func main() {
    emails := []string{
        "test@example.com",
        "duplicate@example.com",
        "test@example.com",
        "invalid-email",
        "another@test.org",
        "duplicate@example.com",
    }

    fmt.Println("Original list:", emails)
    uniqueEmails := DeduplicateStrings(emails)
    fmt.Println("Deduplicated list:", uniqueEmails)

    fmt.Println("\nEmail validation:")
    for _, email := range uniqueEmails {
        valid := ValidateEmail(email)
        fmt.Printf("%s: %v\n", email, valid)
    }
}