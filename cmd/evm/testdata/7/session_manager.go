
package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	Data      map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Manager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		duration: sessionDuration,
	}
}

func (m *Manager) Create(userID int, initialData map[string]interface{}) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := &Session{
		ID:        token,
		UserID:    userID,
		Data:      initialData,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.duration),
	}

	m.mutex.Lock()
	m.sessions[token] = session
	m.mutex.Unlock()

	return token, nil
}

func (m *Manager) Get(token string) (*Session, error) {
	m.mutex.RLock()
	session, exists := m.sessions[token]
	m.mutex.RUnlock()

	if !exists {
		return nil, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		m.Delete(token)
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (m *Manager) Update(token string, data map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	session, exists := m.sessions[token]
	if !exists {
		return errors.New("session not found")
	}

	for key, value := range data {
		session.Data[key] = value
	}

	return nil
}

func (m *Manager) Delete(token string) {
	m.mutex.Lock()
	delete(m.sessions, token)
	m.mutex.Unlock()
}

func (m *Manager) Cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for token, session := range m.sessions {
		if now.After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}