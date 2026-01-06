package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type Preference struct {
	Theme     string    `json:"theme"`
	Language  string    `json:"language"`
	Timezone  string    `json:"timezone"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserPreferencesManager struct {
	mu          sync.RWMutex
	preferences map[string]Preference
	cacheTTL    time.Duration
	storagePath string
}

func NewUserPreferencesManager(storagePath string, cacheTTL time.Duration) *UserPreferencesManager {
	return &UserPreferencesManager{
		preferences: make(map[string]Preference),
		cacheTTL:    cacheTTL,
		storagePath: storagePath,
	}
}

func (m *UserPreferencesManager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.storagePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &m.preferences)
}

func (m *UserPreferencesManager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.MarshalIndent(m.preferences, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.storagePath, data, 0644)
}

func (m *UserPreferencesManager) Get(userID string) (Preference, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pref, exists := m.preferences[userID]
	if !exists {
		return Preference{}, false
	}

	if time.Since(pref.UpdatedAt) > m.cacheTTL {
		delete(m.preferences, userID)
		return Preference{}, false
	}

	return pref, true
}

func (m *UserPreferencesManager) Set(userID string, pref Preference) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pref.UpdatedAt = time.Now()
	m.preferences[userID] = pref

	return m.Save()
}

func (m *UserPreferencesManager) Delete(userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.preferences, userID)
	return m.Save()
}

func (m *UserPreferencesManager) GetAllUsers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]string, 0, len(m.preferences))
	for userID := range m.preferences {
		users = append(users, userID)
	}
	return users
}