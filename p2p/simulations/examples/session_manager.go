package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidToken    = errors.New("invalid session token")
)

type Session struct {
	UserID    string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Manager struct {
	client    *redis.Client
	prefix    string
	expiry    time.Duration
}

func NewManager(client *redis.Client, prefix string, expiry time.Duration) *Manager {
	return &Manager{
		client: client,
		prefix: prefix,
		expiry: expiry,
	}
}

func (m *Manager) CreateSession(userID, username string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	session := Session{
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.expiry),
	}

	key := m.prefix + token
	ctx := context.Background()
	
	err = m.client.Set(ctx, key, session, m.expiry).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *Manager) GetSession(token string) (*Session, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	key := m.prefix + token
	ctx := context.Background()

	var session Session
	err := m.client.Get(ctx, key).Scan(&session)
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (m *Manager) DeleteSession(token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	key := m.prefix + token
	ctx := context.Background()

	_, err := m.client.Del(ctx, key).Result()
	return err
}

func (m *Manager) RefreshSession(token string) error {
	session, err := m.GetSession(token)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(m.expiry)
	key := m.prefix + token
	ctx := context.Background()

	return m.client.Set(ctx, key, session, m.expiry).Err()
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}package session

import (
    "sync"
    "time"
)

type Session struct {
    ID        string
    Data      map[string]interface{}
    ExpiresAt time.Time
}

type Manager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
    ttl      time.Duration
}

func NewManager(ttl time.Duration) *Manager {
    m := &Manager{
        sessions: make(map[string]*Session),
        ttl:      ttl,
    }
    go m.cleanupLoop()
    return m
}

func (m *Manager) Create(id string) *Session {
    m.mu.Lock()
    defer m.mu.Unlock()

    session := &Session{
        ID:        id,
        Data:      make(map[string]interface{}),
        ExpiresAt: time.Now().Add(m.ttl),
    }
    m.sessions[id] = session
    return session
}

func (m *Manager) Get(id string) *Session {
    m.mu.RLock()
    defer m.mu.RUnlock()

    session, exists := m.sessions[id]
    if !exists || time.Now().After(session.ExpiresAt) {
        return nil
    }
    return session
}

func (m *Manager) cleanupLoop() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        m.mu.Lock()
        now := time.Now()
        for id, session := range m.sessions {
            if now.After(session.ExpiresAt) {
                delete(m.sessions, id)
            }
        }
        m.mu.Unlock()
    }
}