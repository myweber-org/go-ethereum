package main

import (
    "log"
    "time"
    "context"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
    dbConnectionString = "postgresql://user:pass@localhost/sessions?sslmode=disable"
    cleanupInterval    = 24 * time.Hour
    retentionPeriod    = 30 * 24 * time.Hour
)

func cleanupExpiredSessions(db *sql.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    query := `DELETE FROM user_sessions WHERE last_activity < $1`
    cutoffTime := time.Now().Add(-retentionPeriod)

    result, err := db.ExecContext(ctx, query, cutoffTime)
    if err != nil {
        return err
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions older than %v", rowsAffected, cutoffTime)
    return nil
}

func main() {
    db, err := sql.Open("postgres", dbConnectionString)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("Database connection test failed:", err)
    }

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    log.Println("Session cleanup service started")
    for range ticker.C {
        if err := cleanupExpiredSessions(db); err != nil {
            log.Printf("Cleanup failed: %v", err)
        }
    }
}