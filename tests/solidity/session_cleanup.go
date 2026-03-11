
package main

import (
	"context"
	"log"
	"time"

	"your_project/internal/db"
	"your_project/internal/models"
)

func main() {
	ctx := context.Background()
	database, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	for {
		expiredTime := time.Now().Add(-24 * time.Hour)
		result, err := database.ExecContext(ctx,
			"DELETE FROM user_sessions WHERE last_activity < ?",
			expiredTime,
		)
		if err != nil {
			log.Printf("Error cleaning sessions: %v", err)
		} else {
			rows, _ := result.RowsAffected()
			log.Printf("Cleaned %d expired sessions", rows)
		}

		time.Sleep(24 * time.Hour)
	}
}