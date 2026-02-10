package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type UserPreferences struct {
	Theme      string `json:"theme"`
	Language   string `json:"language"`
	Timezone   string `json:"timezone"`
	EmailAlerts bool  `json:"email_alerts"`
}

func LoadPreferences(filename string) (*UserPreferences, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var prefs UserPreferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	if prefs.Theme == "" {
		prefs.Theme = "light"
	}
	if prefs.Language == "" {
		prefs.Language = "en"
	}
	if prefs.Timezone == "" {
		prefs.Timezone = "UTC"
	}

	return &prefs, nil
}

func main() {
	prefs, err := LoadPreferences("preferences.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded preferences:\n")
	fmt.Printf("Theme: %s\n", prefs.Theme)
	fmt.Printf("Language: %s\n", prefs.Language)
	fmt.Printf("Timezone: %s\n", prefs.Timezone)
	fmt.Printf("Email Alerts: %v\n", prefs.EmailAlerts)
}