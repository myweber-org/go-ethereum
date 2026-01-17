package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/database"
	"yourproject/internal/models"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	ctx := context.Background()
	cutoff := time.Now().Add(-24 * time.Hour)

	result, err := db.ExecContext(ctx,
		"DELETE FROM user_sessions WHERE last_activity < ?",
		cutoff)
	if err != nil {
		log.Printf("Session cleanup failed: %v", err)
		return
	}

	rows, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", rows)
}