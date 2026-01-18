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
}