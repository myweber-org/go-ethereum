
package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type DataProcessor struct {
	mu    sync.Mutex
	items []string
}

func (dp *DataProcessor) AddItem(item string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	dp.items = append(dp.items, item)
}

func (dp *DataProcessor) ProcessItems() {
	var wg sync.WaitGroup
	for _, item := range dp.items {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			fmt.Printf("Processing: %s\n", s)
		}(item)
	}
	wg.Wait()
}

func main() {
	processor := &DataProcessor{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter items to process (type 'done' to finish):")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "done" {
			break
		}
		processor.AddItem(line)
	}

	fmt.Println("\nProcessing items concurrently:")
	processor.ProcessItems()
	fmt.Println("All items processed successfully")
}