package main

import (
    "encoding/json"
    "sync"
    "time"
)

type UserPreferences struct {
    UserID    string                 `json:"user_id"`
    Settings  map[string]interface{} `json:"settings"`
    UpdatedAt time.Time              `json:"updated_at"`
}

type PreferenceCache struct {
    mu      sync.RWMutex
    store   map[string]*UserPreferences
    ttl     time.Duration
    cleanup time.Duration
}

func NewPreferenceCache(ttl, cleanupInterval time.Duration) *PreferenceCache {
    cache := &PreferenceCache{
        store:   make(map[string]*UserPreferences),
        ttl:     ttl,
        cleanup: cleanupInterval,
    }
    go cache.startCleanup()
    return cache
}

func (c *PreferenceCache) Set(prefs *UserPreferences) {
    c.mu.Lock()
    defer c.mu.Unlock()
    prefs.UpdatedAt = time.Now()
    c.store[prefs.UserID] = prefs
}

func (c *PreferenceCache) Get(userID string) (*UserPreferences, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    prefs, exists := c.store[userID]
    if !exists {
        return nil, false
    }
    if time.Since(prefs.UpdatedAt) > c.ttl {
        return nil, false
    }
    return prefs, true
}

func (c *PreferenceCache) Delete(userID string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.store, userID)
}

func (c *PreferenceCache) startCleanup() {
    ticker := time.NewTicker(c.cleanup)
    defer ticker.Stop()
    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for userID, prefs := range c.store {
            if now.Sub(prefs.UpdatedAt) > c.ttl {
                delete(c.store, userID)
            }
        }
        c.mu.Unlock()
    }
}

func (c *PreferenceCache) MarshalJSON() ([]byte, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    data := make(map[string]*UserPreferences)
    for k, v := range c.store {
        data[k] = v
    }
    return json.Marshal(data)
}