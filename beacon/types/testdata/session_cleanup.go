package main

import (
    "database/sql"
    "log"
    "time"
)

func cleanupExpiredSessions(db *sql.DB) error {
    query := `DELETE FROM user_sessions WHERE expires_at < $1`
    result, err := db.Exec(query, time.Now())
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    log.Printf("Cleaned up %d expired sessions", rowsAffected)
    return nil
}

func main() {
    db, err := sql.Open("postgres", "user=postgres dbname=app sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := cleanupExpiredSessions(db); err != nil {
        log.Fatal(err)
    }
}