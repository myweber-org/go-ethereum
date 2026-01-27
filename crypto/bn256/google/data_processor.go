package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Data string
}

type Result struct {
	TaskID int
	Output string
	Err    error
}

func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d\n", id, task.ID)
		time.Sleep(100 * time.Millisecond)
		results <- Result{
			TaskID: task.ID,
			Output: fmt.Sprintf("Processed: %s", task.Data),
		}
	}
}

func processTasks(tasks []Task, numWorkers int) []Result {
	taskChan := make(chan Task, len(tasks))
	resultChan := make(chan Result, len(tasks))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskChan, resultChan, &wg)
	}

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []Result
	for result := range resultChan {
		results = append(results, result)
	}
	return results
}

func main() {
	tasks := []Task{
		{ID: 1, Data: "alpha"},
		{ID: 2, Data: "beta"},
		{ID: 3, Data: "gamma"},
		{ID: 4, Data: "delta"},
		{ID: 5, Data: "epsilon"},
	}

	results := processTasks(tasks, 3)
	fmt.Println("\nProcessing complete:")
	for _, r := range results {
		fmt.Printf("Task %d: %s\n", r.TaskID, r.Output)
	}
}