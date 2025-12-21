
package cache

import (
	"sync"
	"time"
)

type UserPreferences struct {
	UserID    string
	Theme     string
	Language  string
	Timezone  string
	UpdatedAt time.Time
}

type PreferencesCache struct {
	mu          sync.RWMutex
	preferences map[string]*UserPreferences
	ttl         time.Duration
	stopChan    chan struct{}
}

func NewPreferencesCache(ttl time.Duration) *PreferencesCache {
	cache := &PreferencesCache{
		preferences: make(map[string]*UserPreferences),
		ttl:         ttl,
		stopChan:    make(chan struct{}),
	}
	go cache.startCleanup()
	return cache
}

func (c *PreferencesCache) Set(userID string, prefs *UserPreferences) {
	c.mu.Lock()
	defer c.mu.Unlock()
	prefs.UpdatedAt = time.Now()
	c.preferences[userID] = prefs
}

func (c *PreferencesCache) Get(userID string) (*UserPreferences, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	prefs, exists := c.preferences[userID]
	if !exists {
		return nil, false
	}
	if time.Since(prefs.UpdatedAt) > c.ttl {
		return nil, false
	}
	return prefs, true
}

func (c *PreferencesCache) Delete(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.preferences, userID)
}

func (c *PreferencesCache) startCleanup() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopChan:
			return
		}
	}
}

func (c *PreferencesCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for userID, prefs := range c.preferences {
		if now.Sub(prefs.UpdatedAt) > c.ttl {
			delete(c.preferences, userID)
		}
	}
}

func (c *PreferencesCache) Stop() {
	close(c.stopChan)
}