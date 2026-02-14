
package main

import (
    "context"
    "database/sql"
    "log"
    "time"

    _ "github.com/lib/pq"
)

func cleanupExpiredSessions(db *sql.DB) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    query := `DELETE FROM user_sessions WHERE expires_at < NOW()`
    result, err := db.ExecContext(ctx, query)
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

func main() {
    connStr := "user=postgres dbname=appdb sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    log.Println("Session cleanup service started")
    for range ticker.C {
        if err := cleanupExpiredSessions(db); err != nil {
            log.Printf("Cleanup error: %v", err)
        }
    }
}