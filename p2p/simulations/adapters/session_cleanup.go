package main

import (
	"context"
	"log"
	"time"
)

type SessionStore interface {
	DeleteExpiredSessions(ctx context.Context) error
}

type CleanupJob struct {
	store     SessionStore
	interval  time.Duration
}

func NewCleanupJob(store SessionStore, interval time.Duration) *CleanupJob {
	return &CleanupJob{
		store:    store,
		interval: interval,
	}
}

func (j *CleanupJob) Run(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Cleanup job stopped")
			return
		case <-ticker.C:
			if err := j.store.DeleteExpiredSessions(ctx); err != nil {
				log.Printf("Failed to delete expired sessions: %v", err)
			} else {
				log.Println("Expired sessions cleaned up successfully")
			}
		}
	}
}

func main() {
	ctx := context.Background()
	store := &mockSessionStore{}
	job := NewCleanupJob(store, 24*time.Hour)

	go job.Run(ctx)

	// Keep main running
	select {}
}

type mockSessionStore struct{}

func (m *mockSessionStore) DeleteExpiredSessions(ctx context.Context) error {
	// Simulate session cleanup
	time.Sleep(100 * time.Millisecond)
	return nil
}
package main

import (
    "context"
    "log"
    "time"
)

type SessionStore interface {
    DeleteExpiredSessions(ctx context.Context) error
}

type SessionCleanupJob struct {
    store     SessionStore
    interval  time.Duration
    logger    *log.Logger
    stopChan  chan struct{}
}

func NewSessionCleanupJob(store SessionStore, interval time.Duration, logger *log.Logger) *SessionCleanupJob {
    return &SessionCleanupJob{
        store:    store,
        interval: interval,
        logger:   logger,
        stopChan: make(chan struct{}),
    }
}

func (j *SessionCleanupJob) Start() {
    go j.run()
}

func (j *SessionCleanupJob) Stop() {
    close(j.stopChan)
}

func (j *SessionCleanupJob) run() {
    ticker := time.NewTicker(j.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            j.executeCleanup()
        case <-j.stopChan:
            j.logger.Println("Session cleanup job stopped")
            return
        }
    }
}

func (j *SessionCleanupJob) executeCleanup() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    start := time.Now()
    err := j.store.DeleteExpiredSessions(ctx)
    duration := time.Since(start)

    if err != nil {
        j.logger.Printf("Failed to clean up expired sessions: %v", err)
    } else {
        j.logger.Printf("Session cleanup completed in %v", duration)
    }
}