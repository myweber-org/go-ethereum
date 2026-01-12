package session

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "sync"
    "time"
)

type Session struct {
    UserID    string
    ExpiresAt time.Time
}

type Manager struct {
    sessions map[string]Session
    mu       sync.RWMutex
    duration time.Duration
}

func NewManager(sessionDuration time.Duration) *Manager {
    return &Manager{
        sessions: make(map[string]Session),
        duration: sessionDuration,
    }
}

func (m *Manager) CreateSession(userID string) (string, error) {
    token, err := generateToken()
    if err != nil {
        return "", err
    }

    m.mu.Lock()
    defer m.mu.Unlock()

    m.sessions[token] = Session{
        UserID:    userID,
        ExpiresAt: time.Now().Add(m.duration),
    }

    return token, nil
}

func (m *Manager) ValidateSession(token string) (string, error) {
    m.mu.RLock()
    session, exists := m.sessions[token]
    m.mu.RUnlock()

    if !exists {
        return "", errors.New("session not found")
    }

    if time.Now().After(session.ExpiresAt) {
        m.mu.Lock()
        delete(m.sessions, token)
        m.mu.Unlock()
        return "", errors.New("session expired")
    }

    return session.UserID, nil
}

func (m *Manager) RemoveSession(token string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.sessions, token)
}

func (m *Manager) CleanupExpired() {
    m.mu.Lock()
    defer m.mu.Unlock()

    now := time.Now()
    for token, session := range m.sessions {
        if now.After(session.ExpiresAt) {
            delete(m.sessions, token)
        }
    }
}

func generateToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}