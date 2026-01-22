package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionCleaner struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionCleaner(addr string, ttl time.Duration) *SessionCleaner {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &SessionCleaner{
		client: rdb,
		ttl:    ttl,
	}
}

func (sc *SessionCleaner) CleanExpiredSessions(ctx context.Context) error {
	iter := sc.client.Scan(ctx, 0, "session:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		ttl, err := sc.client.TTL(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get TTL for key %s: %v", key, err)
			continue
		}
		if ttl <= 0 {
			if err := sc.client.Del(ctx, key).Err(); err != nil {
				log.Printf("Failed to delete expired session %s: %v", key, err)
			} else {
				log.Printf("Cleaned expired session: %s", key)
			}
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func (sc *SessionCleaner) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := sc.CleanExpiredSessions(ctx); err != nil {
				log.Printf("Session cleanup failed: %v", err)
			}
		}
	}
}

func main() {
	cleaner := NewSessionCleaner("localhost:6379", 24*time.Hour)
	ctx := context.Background()
	cleaner.Run(ctx, time.Hour)
}