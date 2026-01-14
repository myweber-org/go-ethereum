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
}