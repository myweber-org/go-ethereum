package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type SessionCleaner struct {
    db        *sql.DB
    interval  time.Duration
    retention time.Duration
}

func NewSessionCleaner(db *sql.DB, interval, retention time.Duration) *SessionCleaner {
    return &SessionCleaner{
        db:        db,
        interval:  interval,
        retention: retention,
    }
}

func (sc *SessionCleaner) Run(ctx context.Context) {
    ticker := time.NewTicker(sc.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Session cleaner stopped")
            return
        case <-ticker.C:
            sc.cleanup()
        }
    }
}

func (sc *SessionCleaner) cleanup() {
    cutoff := time.Now().Add(-sc.retention)
    query := `DELETE FROM user_sessions WHERE last_activity < $1`

    result, err := sc.db.Exec(query, cutoff)
    if err != nil {
        log.Printf("Failed to clean sessions: %v", err)
        return
    }

    rows, _ := result.RowsAffected()
    if rows > 0 {
        log.Printf("Cleaned %d expired sessions", rows)
    }
}package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func cleanupExpiredSessions() error {
	// Simulate finding expired session keys (e.g., via pattern matching)
	// In reality, you might use SCAN with a session key pattern
	iter := rdb.Scan(ctx, 0, "session:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		ttl, err := rdb.TTL(ctx, key).Result()
		if err != nil {
			log.Printf("Error getting TTL for key %s: %v", key, err)
			continue
		}
		if ttl < 0 {
			// Key has no expiration or is persistent, skip
			continue
		}
		if ttl == 0 {
			// Key has expired, delete it
			if err := rdb.Del(ctx, key).Err(); err != nil {
				log.Printf("Error deleting expired key %s: %v", key, err)
			} else {
				log.Printf("Deleted expired session: %s", key)
			}
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	initRedis()
	defer rdb.Close()

	// Run cleanup every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Starting session cleanup...")
			if err := cleanupExpiredSessions(); err != nil {
				log.Printf("Cleanup failed: %v", err)
			} else {
				log.Println("Cleanup completed successfully")
			}
		}
	}
}