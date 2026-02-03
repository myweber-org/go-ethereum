package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type UserDataAggregator struct {
	users []User
	mu    sync.RWMutex
}

func NewUserDataAggregator() *UserDataAggregator {
	return &UserDataAggregator{
		users: make([]User, 0),
	}
}

func (uda *UserDataAggregator) AddUser(user User) {
	uda.mu.Lock()
	defer uda.mu.Unlock()
	uda.users = append(uda.users, user)
}

func (uda *UserDataAggregator) GetUsers() []User {
	uda.mu.RLock()
	defer uda.mu.RUnlock()
	return uda.users
}

func (uda *UserDataAggregator) AggregateUserData(userIDs []int) map[int]User {
	result := make(map[int]User)
	var wg sync.WaitGroup
	userChan := make(chan User, len(userIDs))

	for _, id := range userIDs {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			if user := uda.fetchUserData(userID); user.ID != 0 {
				userChan <- user
			}
		}(id)
	}

	wg.Wait()
	close(userChan)

	for user := range userChan {
		result[user.ID] = user
	}

	return result
}

func (uda *UserDataAggregator) fetchUserData(userID int) User {
	time.Sleep(50 * time.Millisecond)
	uda.mu.RLock()
	defer uda.mu.RUnlock()

	for _, user := range uda.users {
		if user.ID == userID {
			return user
		}
	}
	return User{}
}

func (uda *UserDataAggregator) ExportJSON() (string, error) {
	uda.mu.RLock()
	defer uda.mu.RUnlock()

	data, err := json.MarshalIndent(uda.users, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	aggregator := NewUserDataAggregator()

	aggregator.AddUser(User{ID: 1, Name: "John Doe", Email: "john@example.com", CreatedAt: time.Now().Format(time.RFC3339)})
	aggregator.AddUser(User{ID: 2, Name: "Jane Smith", Email: "jane@example.com", CreatedAt: time.Now().Format(time.RFC3339)})
	aggregator.AddUser(User{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", CreatedAt: time.Now().Format(time.RFC3339)})

	userData := aggregator.AggregateUserData([]int{1, 2, 3})
	fmt.Printf("Aggregated %d users\n", len(userData))

	jsonData, err := aggregator.ExportJSON()
	if err != nil {
		fmt.Printf("Export error: %v\n", err)
		return
	}

	fmt.Println("Exported user data:")
	fmt.Println(jsonData)
}