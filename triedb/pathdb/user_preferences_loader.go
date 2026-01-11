
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type UserPreferences struct {
	Theme     string `json:"theme"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	NotificationsEnabled bool `json:"notifications_enabled"`
}

type PreferencesCache struct {
	mu    sync.RWMutex
	store map[string]cachedPreference
}

type cachedPreference struct {
	prefs   UserPreferences
	expires time.Time
}

func NewPreferencesCache() *PreferencesCache {
	return &PreferencesCache{
		store: make(map[string]cachedPreference),
	}
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.store[userID]
	if !exists || time.Now().After(item.expires) {
		return UserPreferences{}, false
	}
	return item.prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.store[userID] = cachedPreference{
		prefs:   prefs,
		expires: time.Now().Add(5 * time.Minute),
	}
}

func LoadUserPreferences(userID string, cache *PreferencesCache) (UserPreferences, error) {
	if cached, found := cache.Get(userID); found {
		fmt.Printf("Cache hit for user %s\n", userID)
		return cached, nil
	}

	fmt.Printf("Loading preferences for user %s\n", userID)
	
	prefs, err := fetchPreferencesFromSource(userID)
	if err != nil {
		return UserPreferences{}, fmt.Errorf("failed to load preferences: %w", err)
	}

	if err := validatePreferences(prefs); err != nil {
		return UserPreferences{}, fmt.Errorf("invalid preferences: %w", err)
	}

	cache.Set(userID, prefs)
	return prefs, nil
}

func fetchPreferencesFromSource(userID string) (UserPreferences, error) {
	time.Sleep(100 * time.Millisecond)
	
	prefsJSON := `{"theme":"dark","language":"en","timezone":"UTC","notifications_enabled":true}`
	
	var prefs UserPreferences
	if err := json.Unmarshal([]byte(prefsJSON), &prefs); err != nil {
		return UserPreferences{}, err
	}
	
	return prefs, nil
}

func validatePreferences(prefs UserPreferences) error {
	if prefs.Theme == "" {
		return errors.New("theme cannot be empty")
	}
	if prefs.Language == "" {
		return errors.New("language cannot be empty")
	}
	if prefs.Timezone == "" {
		return errors.New("timezone cannot be empty")
	}
	return nil
}

func main() {
	cache := NewPreferencesCache()
	
	prefs, err := LoadUserPreferences("user123", cache)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Loaded preferences: %+v\n", prefs)
	
	prefs2, _ := LoadUserPreferences("user123", cache)
	fmt.Printf("Second load (cached): %+v\n", prefs2)
}