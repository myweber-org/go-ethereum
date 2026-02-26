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

func (d *DatabaseSessionStore) DeleteExpiredSessions() error {
	// Simulate database cleanup operation
	log.Println("Deleting expired sessions from database")
	return nil
}

func scheduleSessionCleanup(store SessionStore, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := store.DeleteExpiredSessions(); err != nil {
			log.Printf("Failed to clean up sessions: %v", err)
		} else {
			log.Println("Session cleanup completed successfully")
		}
	}
}

func main() {
	store := &DatabaseSessionStore{}
	scheduleSessionCleanup(store, 24*time.Hour)
}