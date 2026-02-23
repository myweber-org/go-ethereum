package main

import (
	"fmt"
	"regexp"
	"strings"
)

func SanitizeFilename(filename string) string {
	// Remove any characters not alphanumeric, dash, underscore, or dot
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-_.]`)
	sanitized := reg.ReplaceAllString(filename, "")

	// Replace multiple dots or dashes with a single one
	regMulti := regexp.MustCompile(`\.{2,}`)
	sanitized = regMulti.ReplaceAllString(sanitized, ".")

	regMultiDash := regexp.MustCompile(`\-{2,}`)
	sanitized = regMultiDash.ReplaceAllString(sanitized, "-")

	// Trim leading/trailing dots and dashes
	sanitized = strings.Trim(sanitized, ".-_")

	// If empty after sanitization, return default name
	if sanitized == "" {
		return "untitled"
	}

	return sanitized
}

func main() {
	testCases := []string{
		"my file*.txt",
		"---test---.pdf",
		"con...fig.yml",
		"<script>alert()</script>.exe",
		"normal_name-123.jpg",
		"",
	}

	for _, tc := range testCases {
		result := SanitizeFilename(tc)
		fmt.Printf("Original: %-30s -> Sanitized: %s\n", tc, result)
	}
}