package main

import (
    "context"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
    rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
}

func cleanupExpiredSessions() error {
    // Assume sessions are stored with prefix "session:" and have TTL
    // This is a simple scan-based cleanup for demonstration
    // In production, consider using Redis keyspace notifications or other methods
    iter := rdb.Scan(ctx, 0, "session:*", 0).Iterator()
    for iter.Next(ctx) {
        key := iter.Val()
        ttl, err := rdb.TTL(ctx, key).Result()
        if err != nil {
            log.Printf("Failed to get TTL for key %s: %v", key, err)
            continue
        }
        if ttl < 0 {
            // Key has no expiration, remove it
            if err := rdb.Del(ctx, key).Err(); err != nil {
                log.Printf("Failed to delete key %s: %v", key, err)
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
    // Run cleanup daily at 3 AM
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            log.Println("Starting session cleanup...")
            if err := cleanupExpiredSessions(); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            } else {
                log.Println("Session cleanup completed")
            }
        }
    }
}