package main

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type UserPreferences struct {
	Theme     string `json:"theme"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
	UpdatedAt int64  `json:"updated_at"`
}

type PreferencesCache struct {
	mu      sync.RWMutex
	entries map[string]UserPreferences
	ttl     time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	return &PreferencesCache{
		entries: make(map[string]UserPreferences),
		ttl:     ttl,
	}
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
	c.mu.RLock()
	prefs, exists := c.entries[userID]
	c.mu.RUnlock()

	if !exists {
		return UserPreferences{}, false
	}

	if time.Now().Unix()-prefs.UpdatedAt > int64(c.ttl.Seconds()) {
		c.mu.Lock()
		delete(c.entries, userID)
		c.mu.Unlock()
		return UserPreferences{}, false
	}

	return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
	prefs.UpdatedAt = time.Now().Unix()
	c.mu.Lock()
	c.entries[userID] = prefs
	c.mu.Unlock()
}

func LoadPreferences(userID string, cache *PreferencesCache) (UserPreferences, error) {
	if cache != nil {
		if prefs, found := cache.Get(userID); found {
			return prefs, nil
		}
	}

	prefs, err := fetchPreferencesFromStorage(userID)
	if err != nil {
		return UserPreferences{}, err
	}

	if err := validatePreferences(prefs); err != nil {
		return UserPreferences{}, err
	}

	if cache != nil {
		cache.Set(userID, prefs)
	}

	return prefs, nil
}

func fetchPreferencesFromStorage(userID string) (UserPreferences, error) {
	if userID == "" {
		return UserPreferences{}, errors.New("invalid user identifier")
	}

	return UserPreferences{
		Theme:     "dark",
		Language:  "en",
		Timezone:  "UTC",
		UpdatedAt: time.Now().Unix(),
	}, nil
}

func validatePreferences(prefs UserPreferences) error {
	if prefs.Theme != "light" && prefs.Theme != "dark" {
		return errors.New("invalid theme selection")
	}
	if prefs.Language != "en" && prefs.Language != "es" && prefs.Language != "fr" {
		return errors.New("unsupported language")
	}
	return nil
}

func main() {
	cache := NewPreferencesCache(5 * time.Minute)
	prefs, err := LoadPreferences("user123", cache)
	if err != nil {
		panic(err)
	}

	jsonData, _ := json.MarshalIndent(prefs, "", "  ")
	println(string(jsonData))
}