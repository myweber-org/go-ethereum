package main

import (
    "log"
    "time"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
    cleanupInterval = 24 * time.Hour
    sessionTTL      = 7 * 24 * time.Hour
    deleteBatchSize = 1000
)

func cleanupExpiredSessions(db *sql.DB) error {
    cutoffTime := time.Now().Add(-sessionTTL)
    
    for {
        result, err := db.Exec(`
            DELETE FROM user_sessions 
            WHERE last_activity < $1 
            AND session_id IN (
                SELECT session_id 
                FROM user_sessions 
                WHERE last_activity < $1 
                LIMIT $2
            )`,
            cutoffTime, deleteBatchSize)
        
        if err != nil {
            return err
        }
        
        rowsAffected, _ := result.RowsAffected()
        if rowsAffected == 0 {
            break
        }
        
        log.Printf("Cleaned up %d expired sessions", rowsAffected)
        time.Sleep(100 * time.Millisecond)
    }
    
    return nil
}

func startSessionCleanup(db *sql.DB) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := cleanupExpiredSessions(db); err != nil {
                log.Printf("Session cleanup failed: %v", err)
            } else {
                log.Println("Session cleanup completed successfully")
            }
        }
    }
}

func main() {
    db, err := sql.Open("postgres", "host=localhost port=5432 user=app dbname=appdb sslmode=disable")
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }
    defer db.Close()
    
    if err := db.Ping(); err != nil {
        log.Fatal("Database ping failed:", err)
    }
    
    log.Println("Starting session cleanup service...")
    startSessionCleanup(db)
}