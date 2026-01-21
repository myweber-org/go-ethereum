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

func (s *CleanupService) RunDailyCleanup(ctx context.Context) error {
    log.Println("Starting daily session cleanup")
    
    if err := s.store.DeleteExpired(ctx); err != nil {
        return err
    }
    
    log.Println("Session cleanup completed successfully")
    return nil
}

func (s *CleanupService) StartScheduler(ctx context.Context) {
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            log.Println("Cleanup scheduler stopped")
            return
        case <-ticker.C:
            if err := s.RunDailyCleanup(ctx); err != nil {
                log.Printf("Cleanup failed: %v", err)
            }
        }
    }
}