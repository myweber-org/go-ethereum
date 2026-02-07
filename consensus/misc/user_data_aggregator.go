package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}

type UserStats struct {
	TotalUsers   int `json:"total_users"`
	ActiveUsers  int `json:"active_users"`
	InactiveUsers int `json:"inactive_users"`
}

func fetchUserData(userID int) (User, error) {
	time.Sleep(50 * time.Millisecond)
	
	if userID%10 == 0 {
		return User{}, fmt.Errorf("failed to fetch user %d", userID)
	}
	
	return User{
		ID:        userID,
		Name:      fmt.Sprintf("User%d", userID),
		Email:     fmt.Sprintf("user%d@example.com", userID),
		Active:    userID%3 != 0,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func processUserBatch(startID, batchSize int, results chan<- User, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for i := 0; i < batchSize; i++ {
		userID := startID + i
		user, err := fetchUserData(userID)
		if err != nil {
			errors <- err
			continue
		}
		results <- user
	}
}

func aggregateUserStats(users []User) UserStats {
	stats := UserStats{}
	for _, user := range users {
		stats.TotalUsers++
		if user.Active {
			stats.ActiveUsers++
		} else {
			stats.InactiveUsers++
		}
	}
	return stats
}

func main() {
	startTime := time.Now()
	
	const totalUsers = 100
	const batchSize = 20
	const workerCount = 5
	
	userChan := make(chan User, totalUsers)
	errorChan := make(chan error, totalUsers)
	var wg sync.WaitGroup
	
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		startID := i*batchSize + 1
		go processUserBatch(startID, batchSize, userChan, errorChan, &wg)
	}
	
	wg.Wait()
	close(userChan)
	close(errorChan)
	
	var users []User
	for user := range userChan {
		users = append(users, user)
	}
	
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}
	
	stats := aggregateUserStats(users)
	
	statsJSON, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Println("User Statistics:")
	fmt.Println(string(statsJSON))
	
	fmt.Printf("\nProcessed %d users with %d errors\n", len(users), len(errors))
	fmt.Printf("Execution time: %v\n", time.Since(startTime))
	
	if len(errors) > 0 {
		log.Printf("Encountered %d errors during processing", len(errors))
	}
}