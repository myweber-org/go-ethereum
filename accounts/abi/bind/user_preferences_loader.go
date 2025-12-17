package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type UserPreferences struct {
	Theme      string   `json:"theme"`
	Language   string   `json:"language"`
	Timezone   string   `json:"timezone"`
	EmailAlerts bool    `json:"email_alerts"`
	Dashboard  []string `json:"dashboard_widgets"`
}

type PreferencesCache struct {
	mu          sync.RWMutex
	preferences map[string]UserPreferences
	lastUpdated map[string]time.Time
	ttl         time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		preferences: make(map[string]UserPreferences),
		lastUpdated: make(map[string]time.Time),
		ttl:         ttl,
	}
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
	c.mu.RLock()
	prefs, exists := c.preferences[userID]
	lastUpdate, timeExists := c.lastUpdated[userID]
	c.mu.RUnlock()

	if !exists || !timeExists || time.Since(lastUpdate) > c.ttl {
		return UserPreferences{}, false
	}
	return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
	c.mu.Lock()
	c.preferences[userID] = prefs
	c.lastUpdated[userID] = time.Now()
	c.mu.Unlock()
}

func LoadPreferencesFromFile(filename string) (UserPreferences, error) {
	var prefs UserPreferences

	file, err := os.Open(filename)
	if err != nil {
		return prefs, fmt.Errorf("failed to open preferences file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&prefs); err != nil {
		return prefs, fmt.Errorf("failed to decode preferences: %w", err)
	}

	if err := ValidatePreferences(prefs); err != nil {
		return prefs, fmt.Errorf("invalid preferences: %w", err)
	}

	return prefs, nil
}

func ValidatePreferences(prefs UserPreferences) error {
	if prefs.Theme == "" {
		return fmt.Errorf("theme cannot be empty")
	}
	if prefs.Language == "" {
		return fmt.Errorf("language cannot be empty")
	}
	if prefs.Timezone == "" {
		return fmt.Errorf("timezone cannot be empty")
	}
	return nil
}

func main() {
	cache := NewPreferencesCache(5 * time.Minute)

	prefs, err := LoadPreferencesFromFile("user_prefs.json")
	if err != nil {
		fmt.Printf("Error loading preferences: %v\n", err)
		return
	}

	cache.Set("user123", prefs)

	if cachedPrefs, found := cache.Get("user123"); found {
		fmt.Printf("Loaded preferences: %+v\n", cachedPrefs)
	} else {
		fmt.Println("Preferences not found in cache")
	}
}