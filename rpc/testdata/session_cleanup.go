package main

import (
    "log"
    "time"
)

type Session struct {
    ID        string
    UserID    int
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

func startCleanupJob(store *SessionStore) {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            store.CleanExpiredSessions()
        }
    }
}

func main() {
    store := NewSessionStore()
    go startCleanupJob(store)

    select {}
}