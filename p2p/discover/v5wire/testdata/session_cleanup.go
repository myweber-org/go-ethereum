package main

import (
    "context"
    "database/sql"
    "log"
    "time"
)

type SessionCleaner struct {
    db        *sql.DB
    batchSize int
    interval  time.Duration
}

func NewSessionCleaner(db *sql.DB, batchSize int, interval time.Duration) *SessionCleaner {
    return &SessionCleaner{
        db:        db,
        batchSize: batchSize,
        interval:  interval,
    }
}

func (sc *SessionCleaner) Start(ctx context.Context) {
    ticker := time.NewTicker(sc.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Session cleaner stopped")
            return
        case <-ticker.C:
            sc.cleanupExpiredSessions()
        }
    }
}

func (sc *SessionCleaner) cleanupExpiredSessions() {
    query := `
        DELETE FROM user_sessions 
        WHERE expires_at < NOW() 
        AND session_id IN (
            SELECT session_id 
            FROM user_sessions 
            WHERE expires_at < NOW() 
            LIMIT $1
        )
    `

    result, err := sc.db.Exec(query, sc.batchSize)
    if err != nil {
        log.Printf("Failed to clean expired sessions: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected > 0 {
        log.Printf("Cleaned %d expired sessions", rowsAffected)
    }
}

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/sessions")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    cleaner := NewSessionCleaner(db, 1000, 5*time.Minute)
    ctx := context.Background()
    cleaner.Start(ctx)
}