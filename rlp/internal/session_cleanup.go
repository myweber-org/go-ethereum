package session

import (
    "context"
    "database/sql"
    "log"
    "time"
)

const cleanupInterval = 24 * time.Hour

type Cleaner struct {
    db *sql.DB
}

func NewCleaner(db *sql.DB) *Cleaner {
    return &Cleaner{db: db}
}

func (c *Cleaner) Start(ctx context.Context) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Session cleanup stopped")
            return
        case <-ticker.C:
            c.cleanupExpiredSessions()
        }
    }
}

func (c *Cleaner) cleanupExpiredSessions() {
    query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
    result, err := c.db.Exec(query)
    if err != nil {
        log.Printf("Failed to clean expired sessions: %v", err)
        return
    }

    rows, _ := result.RowsAffected()
    if rows > 0 {
        log.Printf("Cleaned %d expired sessions", rows)
    }
}