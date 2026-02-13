package main

import (
    "context"
    "log"
    "time"
)

type SessionStore interface {
    DeleteExpired(ctx context.Context) error
}

type CleanupService struct {
    store SessionStore
}

func NewCleanupService(store SessionStore) *CleanupService {
    return &CleanupService{store: store}
}

func (s *CleanupService) RunDailyCleanup() {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            err := s.store.DeleteExpired(ctx)
            cancel()

            if err != nil {
                log.Printf("Failed to delete expired sessions: %v", err)
            } else {
                log.Println("Successfully cleaned up expired sessions")
            }
        }
    }
}

func main() {
    // Implementation would inject actual session store
    // store := NewRedisSessionStore()
    // service := NewCleanupService(store)
    // service.RunDailyCleanup()
}