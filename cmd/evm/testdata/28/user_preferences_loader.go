package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type UserPreferences struct {
	Theme     string `json:"theme"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	UpdatedAt int64  `json:"updated_at"`
}

type PreferencesCache struct {
	preferences map[string]UserPreferences
	lastUpdated time.Time
	ttl         time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		preferences: make(map[string]UserPreferences),
		ttl:         ttl,
	}
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
	prefs, exists := c.preferences[userID]
	if !exists {
		return UserPreferences{}, false
	}
	if time.Since(c.lastUpdated) > c.ttl {
		delete(c.preferences, userID)
		return UserPreferences{}, false
	}
	return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
	c.preferences[userID] = prefs
	c.lastUpdated = time.Now()
}

func LoadPreferencesFromFile(filename string) (map[string]UserPreferences, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var preferences map[string]UserPreferences
	if err := json.Unmarshal(data, &preferences); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for userID, prefs := range preferences {
		if prefs.Theme == "" {
			prefs.Theme = "light"
		}
		if prefs.Language == "" {
			prefs.Language = "en"
		}
		if prefs.Timezone == "" {
			prefs.Timezone = "UTC"
		}
		if prefs.UpdatedAt == 0 {
			prefs.UpdatedAt = time.Now().Unix()
		}
		preferences[userID] = prefs
	}

	return preferences, nil
}

func ValidatePreferences(prefs UserPreferences) error {
	validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
	if !validThemes[prefs.Theme] {
		return fmt.Errorf("invalid theme: %s", prefs.Theme)
	}

	validLanguages := map[string]bool{"en": true, "es": true, "fr": true, "de": true}
	if !validLanguages[prefs.Language] {
		return fmt.Errorf("invalid language: %s", prefs.Language)
	}

	if prefs.UpdatedAt > time.Now().Unix() {
		return fmt.Errorf("invalid update timestamp")
	}

	return nil
}

func main() {
	cache := NewPreferencesCache(5 * time.Minute)

	prefs, err := LoadPreferencesFromFile("preferences.json")
	if err != nil {
		fmt.Printf("Error loading preferences: %v\n", err)
		return
	}

	for userID, userPrefs := range prefs {
		if err := ValidatePreferences(userPrefs); err != nil {
			fmt.Printf("Invalid preferences for user %s: %v\n", userID, err)
			continue
		}
		cache.Set(userID, userPrefs)
		fmt.Printf("Loaded preferences for user %s: %+v\n", userID, userPrefs)
	}

	if cachedPrefs, found := cache.Get("user123"); found {
		fmt.Printf("Retrieved from cache: %+v\n", cachedPrefs)
	}
}