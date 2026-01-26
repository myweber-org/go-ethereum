package main

import (
    "encoding/json"
    "errors"
    "sync"
    "time"
)

type UserPreference struct {
    UserID      string                 `json:"user_id"`
    Theme       string                 `json:"theme"`
    Language    string                 `json:"language"`
    NotificationsEnabled bool          `json:"notifications_enabled"`
    CustomSettings map[string]any      `json:"custom_settings"`
    LastUpdated time.Time              `json:"last_updated"`
}

type PreferenceManager struct {
    cache      map[string]UserPreference
    cacheMutex sync.RWMutex
    cacheTTL   time.Duration
}

func NewPreferenceManager(ttl time.Duration) *PreferenceManager {
    return &PreferenceManager{
        cache:    make(map[string]UserPreference),
        cacheTTL: ttl,
    }
}

func (pm *PreferenceManager) ValidatePreference(pref UserPreference) error {
    if pref.UserID == "" {
        return errors.New("user ID cannot be empty")
    }
    
    validThemes := map[string]bool{"light": true, "dark": true, "auto": true}
    if !validThemes[pref.Theme] {
        return errors.New("invalid theme selection")
    }
    
    validLanguages := map[string]bool{"en": true, "es": true, "fr": true, "de": true}
    if !validLanguages[pref.Language] {
        return errors.New("unsupported language")
    }
    
    return nil
}

func (pm *PreferenceManager) SetPreference(pref UserPreference) error {
    if err := pm.ValidatePreference(pref); err != nil {
        return err
    }
    
    pref.LastUpdated = time.Now()
    
    pm.cacheMutex.Lock()
    pm.cache[pref.UserID] = pref
    pm.cacheMutex.Unlock()
    
    return nil
}

func (pm *PreferenceManager) GetPreference(userID string) (UserPreference, bool) {
    pm.cacheMutex.RLock()
    pref, exists := pm.cache[userID]
    pm.cacheMutex.RUnlock()
    
    if exists && time.Since(pref.LastUpdated) > pm.cacheTTL {
        pm.cacheMutex.Lock()
        delete(pm.cache, userID)
        pm.cacheMutex.Unlock()
        return UserPreference{}, false
    }
    
    return pref, exists
}

func (pm *PreferenceManager) ExportPreferences(userID string) ([]byte, error) {
    pref, exists := pm.GetPreference(userID)
    if !exists {
        return nil, errors.New("user preferences not found")
    }
    
    return json.MarshalIndent(pref, "", "  ")
}

func (pm *PreferenceManager) ClearCache() {
    pm.cacheMutex.Lock()
    pm.cache = make(map[string]UserPreference)
    pm.cacheMutex.Unlock()
}

func (pm *PreferenceManager) CacheStats() (int, int) {
    pm.cacheMutex.RLock()
    defer pm.cacheMutex.RUnlock()
    
    total := len(pm.cache)
    expired := 0
    
    for _, pref := range pm.cache {
        if time.Since(pref.LastUpdated) > pm.cacheTTL {
            expired++
        }
    }
    
    return total, expired
}