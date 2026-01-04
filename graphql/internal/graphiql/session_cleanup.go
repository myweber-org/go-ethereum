package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type SessionCleaner struct {
    db *sql.DB
}

func NewSessionCleaner(db *sql.DB) *SessionCleaner {
    return &SessionCleaner{db: db}
}

func (sc *SessionCleaner) CleanupExpiredSessions(ctx context.Context) error {
    query := `DELETE FROM user_sessions WHERE expires_at < $1`
    result, err := sc.db.ExecContext(ctx, query, time.Now().UTC())
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Failed to get rows affected: %v", err)
    } else {
        log.Printf("Cleaned up %d expired sessions", rowsAffected)
    }

    return nil
}

func (sc *SessionCleaner) RunPeriodicCleanup(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := sc.CleanupExpiredSessions(ctx); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            }
        case <-ctx.Done():
            return
        }
    }
}