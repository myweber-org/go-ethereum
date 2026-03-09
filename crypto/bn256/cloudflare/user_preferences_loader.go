package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "sync"
    "time"
)

type UserPreferences struct {
    Theme      string `json:"theme"`
    Language   string `json:"language"`
    NotificationsEnabled bool `json:"notifications_enabled"`
    Timezone   string `json:"timezone"`
}

type PreferencesCache struct {
    mu      sync.RWMutex
    data    map[string]UserPreferences
    ttl     time.Duration
    created map[string]time.Time
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
    return &PreferencesCache{
        data:    make(map[string]UserPreferences),
        created: make(map[string]time.Time),
        ttl:     ttl,
    }
}

func (c *PreferencesCache) Get(userID string) (UserPreferences, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    prefs, exists := c.data[userID]
    if !exists {
        return UserPreferences{}, false
    }

    if time.Since(c.created[userID]) > c.ttl {
        return UserPreferences{}, false
    }

    return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[userID] = prefs
    c.created[userID] = time.Now()
}

func ValidatePreferences(prefs UserPreferences) error {
    validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
    if !validThemes[prefs.Theme] {
        return errors.New("invalid theme selection")
    }

    if prefs.Language == "" {
        return errors.New("language cannot be empty")
    }

    if prefs.Timezone == "" {
        return errors.New("timezone cannot be empty")
    }

    return nil
}

func LoadPreferencesFromJSON(data []byte) (UserPreferences, error) {
    var prefs UserPreferences
    if err := json.Unmarshal(data, &prefs); err != nil {
        return UserPreferences{}, fmt.Errorf("failed to parse preferences: %w", err)
    }

    if err := ValidatePreferences(prefs); err != nil {
        return UserPreferences{}, fmt.Errorf("invalid preferences: %w", err)
    }

    return prefs, nil
}

func LoadUserPreferences(userID string, cache *PreferencesCache, fetchFunc func(string) ([]byte, error)) (UserPreferences, error) {
    if cached, found := cache.Get(userID); found {
        return cached, nil
    }

    data, err := fetchFunc(userID)
    if err != nil {
        return UserPreferences{}, fmt.Errorf("failed to fetch preferences: %w", err)
    }

    prefs, err := LoadPreferencesFromJSON(data)
    if err != nil {
        return UserPreferences{}, err
    }

    cache.Set(userID, prefs)
    return prefs, nil
}

func main() {
    cache := NewPreferencesCache(5 * time.Minute)

    mockFetch := func(userID string) ([]byte, error) {
        prefs := UserPreferences{
            Theme:                 "dark",
            Language:             "en-US",
            NotificationsEnabled: true,
            Timezone:            "America/New_York",
        }
        return json.Marshal(prefs)
    }

    prefs, err := LoadUserPreferences("user123", cache, mockFetch)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Loaded preferences: %+v\n", prefs)
}