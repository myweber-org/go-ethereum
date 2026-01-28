package utils

import "strings"

func CleanInput(input string) string {
    trimmed := strings.TrimSpace(input)
    cleaned := strings.Join(strings.Fields(trimmed), " ")
    return cleaned
}