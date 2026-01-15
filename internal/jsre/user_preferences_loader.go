package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type UserPreferences struct {
	Theme       string `json:"theme"`
	Language    string `json:"language"`
	ItemsPerPage int   `json:"items_per_page"`
	NotificationsEnabled bool `json:"notifications_enabled"`
}

func LoadPreferences(filename string) (*UserPreferences, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var prefs UserPreferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	prefs = applyDefaults(prefs)
	return &prefs, nil
}

func applyDefaults(prefs UserPreferences) UserPreferences {
	if prefs.Theme == "" {
		prefs.Theme = "light"
	}
	if prefs.Language == "" {
		prefs.Language = "en"
	}
	if prefs.ItemsPerPage <= 0 {
		prefs.ItemsPerPage = 25
	}
	return prefs
}

func main() {
	prefs, err := LoadPreferences("config.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Theme: %s\n", prefs.Theme)
	fmt.Printf("Language: %s\n", prefs.Language)
	fmt.Printf("Items per page: %d\n", prefs.ItemsPerPage)
	fmt.Printf("Notifications enabled: %v\n", prefs.NotificationsEnabled)
}