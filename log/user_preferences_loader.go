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
    NotificationsEnabled bool `json:"notifications_enabled"`
}

type PreferencesCache struct {
    mu      sync.RWMutex
    data    map[string]UserPreferences
    expires map[string]time.Time
    ttl     time.Duration
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
    return &PreferencesCache{
        data:    make(map[string]UserPreferences),
        expires: make(map[string]time.Time),
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

    if expiry, ok := c.expires[userID]; ok && time.Now().After(expiry) {
        return UserPreferences{}, false
    }

    return prefs, true
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[userID] = prefs
    c.expires[userID] = time.Now().Add(c.ttl)
}

func ValidatePreferences(prefs UserPreferences) error {
    validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
    if !validThemes[prefs.Theme] {
        return errors.New("invalid theme selection")
    }

    validLanguages := map[string]bool{"en": true, "es": true, "fr": true, "de": true}
    if !validLanguages[prefs.Language] {
        return errors.New("unsupported language")
    }

    if prefs.Timezone == "" {
        return errors.New("timezone cannot be empty")
    }

    return nil
}

func LoadPreferencesFromJSON(data []byte) (UserPreferences, error) {
    var prefs UserPreferences
    if err := json.Unmarshal(data, &prefs); err != nil {
        return UserPreferences{}, err
    }

    if err := ValidatePreferences(prefs); err != nil {
        return UserPreferences{}, err
    }

    return prefs, nil
}

func LoadUserPreferences(userID string, cache *PreferencesCache, fetchFunc func(string) ([]byte, error)) (UserPreferences, error) {
    if cachedPrefs, found := cache.Get(userID); found {
        return cachedPrefs, nil
    }

    data, err := fetchFunc(userID)
    if err != nil {
        return UserPreferences{}, err
    }

    prefs, err := LoadPreferencesFromJSON(data)
    if err != nil {
        return UserPreferences{}, err
    }

    cache.Set(userID, prefs)
    return prefs, nil
}