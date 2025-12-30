package main

import (
    "context"
    "log"
    "time"

    "github.com/yourproject/db"
)

func main() {
    ctx := context.Background()
    database, err := db.NewConnection()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()

    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions(ctx, database)
        }
    }
}

func cleanupExpiredSessions(ctx context.Context, db *db.Connection) {
    query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
    result, err := db.ExecContext(ctx, query)
    if err != nil {
        log.Printf("Failed to clean up sessions: %v", err)
        return
    }

    rows, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rows)
}