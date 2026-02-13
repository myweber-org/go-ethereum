package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	cutoff := time.Now().Add(-24 * time.Hour)

	result, err := db.ExecContext(ctx,
		"DELETE FROM user_sessions WHERE last_activity < ?",
		cutoff)
	if err != nil {
		log.Printf("Failed to clean sessions: %v", err)
		return
	}

	rows, _ := result.RowsAffected()
	log.Printf("Cleaned %d expired sessions", rows)
}