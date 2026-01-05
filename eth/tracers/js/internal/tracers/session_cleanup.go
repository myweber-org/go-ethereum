
package main

import (
	"context"
	"log"
	"time"

	"yourproject/internal/db"
	"yourproject/internal/models"
)

func main() {
	ctx := context.Background()
	dbClient := db.NewClient()
	defer dbClient.Close()

	// Run cleanup every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanupExpiredSessions(ctx, dbClient)
		}
	}
}

func cleanupExpiredSessions(ctx context.Context, dbClient *db.Client) {
	cutoff := time.Now().Add(-24 * time.Hour)
	result := dbClient.Session.
		Delete().
		Where(models.Session.ExpiresAt.LTE(cutoff)).
		Exec(ctx)

	if result.Error != nil {
		log.Printf("Failed to cleanup sessions: %v", result.Error)
		return
	}

	count, _ := result.RowsAffected()
	log.Printf("Cleaned up %d expired sessions", count)
}