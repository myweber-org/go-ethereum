package session

import (
	"sync"
	"time"
)

type Session struct {
	ID        string
	Data      map[string]interface{}
	ExpiresAt time.Time
}

type Cleaner struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	ttl      time.Duration
	stopChan chan struct{}
}

func NewCleaner(ttl time.Duration) *Cleaner {
	return &Cleaner{
		sessions: make(map[string]*Session),
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}
}

func (c *Cleaner) Add(session *Session) {
	c.mu.Lock()
	defer c.mu.Unlock()
	session.ExpiresAt = time.Now().Add(c.ttl)
	c.sessions[session.ID] = session
}

func (c *Cleaner) Get(id string) (*Session, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	session, exists := c.sessions[id]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil, false
	}
	return session, true
}

func (c *Cleaner) Start() {
	ticker := time.NewTicker(c.ttl / 2)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.cleanup()
			case <-c.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *Cleaner) Stop() {
	close(c.stopChan)
}

func (c *Cleaner) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for id, session := range c.sessions {
		if now.After(session.ExpiresAt) {
			delete(c.sessions, id)
		}
	}
}package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/go-redis/redis/v8"
)

const (
    sessionPrefix = "session:"
    sessionTTL    = 24 * time.Hour
)

func cleanupExpiredSessions(rdb *redis.Client) error {
    ctx := context.Background()
    pattern := sessionPrefix + "*"

    iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
    for iter.Next(ctx) {
        key := iter.Val()
        ttl, err := rdb.TTL(ctx, key).Result()
        if err != nil {
            log.Printf("Failed to get TTL for key %s: %v", key, err)
            continue
        }

        if ttl < 0 {
            if err := rdb.Del(ctx, key).Err(); err != nil {
                log.Printf("Failed to delete expired session %s: %v", key, err)
            } else {
                log.Printf("Removed expired session: %s", key)
            }
        }
    }

    if err := iter.Err(); err != nil {
        return fmt.Errorf("iteration error: %w", err)
    }

    return nil
}

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    defer rdb.Close()

    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        if err := cleanupExpiredSessions(rdb); err != nil {
            log.Printf("Session cleanup failed: %v", err)
        }
    }
}