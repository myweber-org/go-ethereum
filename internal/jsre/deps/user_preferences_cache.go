package cache

import (
	"encoding/json"
	"errors"
	"time"
)

type UserPreferences struct {
	UserID      string                 `json:"user_id"`
	Theme       string                 `json:"theme"`
	Language    string                 `json:"language"`
	Timezone    string                 `json:"timezone"`
	Preferences map[string]interface{} `json:"preferences"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type PreferencesCache struct {
	store    map[string][]byte
	duration time.Duration
}

func NewPreferencesCache(cacheDuration time.Duration) *PreferencesCache {
	return &PreferencesCache{
		store:    make(map[string][]byte),
		duration: cacheDuration,
	}
}

func (c *PreferencesCache) Set(userID string, prefs UserPreferences) error {
	prefs.UpdatedAt = time.Now()
	data, err := json.Marshal(prefs)
	if err != nil {
		return err
	}
	c.store[userID] = data
	return nil
}

func (c *PreferencesCache) Get(userID string) (*UserPreferences, error) {
	data, exists := c.store[userID]
	if !exists {
		return nil, errors.New("preferences not found in cache")
	}

	var prefs UserPreferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, err
	}

	if time.Since(prefs.UpdatedAt) > c.duration {
		delete(c.store, userID)
		return nil, errors.New("cache expired")
	}

	return &prefs, nil
}

func (c *PreferencesCache) Invalidate(userID string) {
	delete(c.store, userID)
}

func (c *PreferencesCache) Clear() {
	c.store = make(map[string][]byte)
}