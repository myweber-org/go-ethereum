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

type SessionStore interface {
	DeleteExpiredSessions() error
}

type DatabaseSessionStore struct{}

func (s *DatabaseSessionStore) DeleteExpiredSessions() error {
	log.Println("Deleting expired sessions from database")
	return nil
}

func cleanupExpiredSessions(store SessionStore) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := store.DeleteExpiredSessions(); err != nil {
			log.Printf("Failed to delete expired sessions: %v", err)
		} else {
			log.Println("Successfully cleaned up expired sessions")
		}
	}
}

func main() {
	store := &DatabaseSessionStore{}
	cleanupExpiredSessions(store)
}