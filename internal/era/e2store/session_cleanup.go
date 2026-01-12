package main

import (
	"log"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type SessionStore struct {
	sessions map[string]Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]Session),
	}
}

func (s *SessionStore) CleanExpiredSessions() {
	now := time.Now()
	count := 0
	for id, session := range s.sessions {
		if session.ExpiresAt.Before(now) {
			delete(s.sessions, id)
			count++
		}
	}
	log.Printf("Cleaned %d expired sessions", count)
}

func main() {
	store := NewSessionStore()
	
	// Simulate adding some sessions
	store.sessions["s1"] = Session{
		ID:        "s1",
		UserID:    "user1",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
	}
	store.sessions["s2"] = Session{
		ID:        "s2",
		UserID:    "user2",
		ExpiresAt: time.Now().Add(1 * time.Hour), // Still valid
	}
	
	// Run cleanup
	store.CleanExpiredSessions()
	
	log.Printf("Remaining sessions: %d", len(store.sessions))
}