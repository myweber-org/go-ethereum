package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

const cleanupInterval = 24 * time.Hour

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredSessions(db)
		}
	}
}

func cleanupExpiredSessions(db *database.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Failed to clean up sessions: %v", err)
		return
	}

	rows, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rows)
}package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	sessionKeyPattern = "session:*"
	batchSize         = 100
	scanInterval      = 1 * time.Hour
)

type SessionCleaner struct {
	client *redis.Client
}

func NewSessionCleaner(addr string) *SessionCleaner {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &SessionCleaner{client: rdb}
}

func (sc *SessionCleaner) CleanExpiredSessions(ctx context.Context) error {
	var cursor uint64
	var deletedCount int

	for {
		keys, nextCursor, err := sc.client.Scan(ctx, cursor, sessionKeyPattern, batchSize).Result()
		if err != nil {
			return fmt.Errorf("scan error: %w", err)
		}

		for _, key := range keys {
			ttl, err := sc.client.TTL(ctx, key).Result()
			if err != nil {
				log.Printf("failed to get TTL for key %s: %v", key, err)
				continue
			}

			if ttl < 0 {
				if err := sc.client.Del(ctx, key).Err(); err != nil {
					log.Printf("failed to delete expired session %s: %v", key, err)
				} else {
					deletedCount++
				}
			}
		}

		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	log.Printf("cleaned %d expired sessions", deletedCount)
	return nil
}

func (sc *SessionCleaner) RunPeriodicCleanup(ctx context.Context) {
	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := sc.CleanExpiredSessions(ctx); err != nil {
				log.Printf("session cleanup failed: %v", err)
			}
		case <-ctx.Done():
			log.Println("session cleanup stopped")
			return
		}
	}
}

func main() {
	ctx := context.Background()
	cleaner := NewSessionCleaner("localhost:6379")

	go cleaner.RunPeriodicCleanup(ctx)

	select {}
}