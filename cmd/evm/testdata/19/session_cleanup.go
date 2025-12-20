package main

import (
    "context"
    "log"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

const (
    dbURL            = "postgresql://user:pass@localhost:5432/appdb"
    cleanupInterval  = 1 * time.Hour
    sessionLifetime  = 24 * time.Hour
    deleteBatchSize  = 1000
)

func main() {
    ctx := context.Background()

    pool, err := pgxpool.New(ctx, dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    defer pool.Close()

    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for range ticker.C {
        err := deleteExpiredSessions(ctx, pool)
        if err != nil {
            log.Printf("Session cleanup failed: %v", err)
        } else {
            log.Println("Session cleanup completed successfully")
        }
    }
}

func deleteExpiredSessions(ctx context.Context, pool *pgxpool.Pool) error {
    cutoffTime := time.Now().Add(-sessionLifetime)

    for {
        result, err := pool.Exec(ctx,
            `DELETE FROM user_sessions 
             WHERE last_activity < $1 
             AND session_id IN (
                 SELECT session_id FROM user_sessions 
                 WHERE last_activity < $1 
                 LIMIT $2
             )`,
            cutoffTime, deleteBatchSize,
        )

        if err != nil {
            return err
        }

        rowsAffected := result.RowsAffected()
        if rowsAffected == 0 {
            break
        }

        log.Printf("Deleted %d expired sessions", rowsAffected)

        if rowsAffected < int64(deleteBatchSize) {
            break
        }
    }

    return nil
}