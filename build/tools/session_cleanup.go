package main

import (
    "database/sql"
    "log"
    "time"
)

var db *sql.DB

func cleanupExpiredSessions() {
    query := `DELETE FROM user_sessions WHERE expires_at < $1`
    result, err := db.Exec(query, time.Now())
    if err != nil {
        log.Printf("Failed to clean up sessions: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    log.Printf("Cleaned up %d expired sessions", rowsAffected)
}

func scheduleCleanup(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            cleanupExpiredSessions()
        }
    }
}

func main() {
    // Database initialization would happen here
    // db = initializeDB()
    
    // Run cleanup every hour
    go scheduleCleanup(time.Hour)
    
    // Keep main goroutine alive
    select {}
}