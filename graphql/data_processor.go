
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DataProcessor struct {
	workers   int
	inputChan chan int
	outputChan chan int
	wg        sync.WaitGroup
	errChan   chan error
}

func NewDataProcessor(workers int) *DataProcessor {
	return &DataProcessor{
		workers:   workers,
		inputChan: make(chan int, 100),
		outputChan: make(chan int, 100),
		errChan:   make(chan error, workers),
	}
}

func (dp *DataProcessor) Start() {
	for i := 0; i < dp.workers; i++ {
		dp.wg.Add(1)
		go dp.worker(i)
	}
}

func (dp *DataProcessor) worker(id int) {
	defer dp.wg.Done()
	for data := range dp.inputChan {
		if data < 0 {
			dp.errChan <- fmt.Errorf("worker %d: negative value %d", id, data)
			continue
		}
		processed := data * 2
		time.Sleep(10 * time.Millisecond)
		dp.outputChan <- processed
	}
}

func (dp *DataProcessor) Process(data []int) ([]int, error) {
	results := make([]int, 0, len(data))
	done := make(chan bool)
	
	go func() {
		for _, val := range data {
			dp.inputChan <- val
		}
		close(dp.inputChan)
	}()
	
	go func() {
		for i := 0; i < len(data); i++ {
			select {
			case result := <-dp.outputChan:
				results = append(results, result)
			case err := <-dp.errChan:
				fmt.Printf("Processing error: %v\n", err)
			}
		}
		done <- true
	}()
	
	<-done
	dp.wg.Wait()
	close(dp.outputChan)
	close(dp.errChan)
	
	if len(results) == 0 {
		return nil, errors.New("no valid results produced")
	}
	
	return results, nil
}

func main() {
	processor := NewDataProcessor(4)
	processor.Start()
	
	data := []int{1, 2, 3, -4, 5, 6, -7, 8}
	
	results, err := processor.Process(data)
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		return
	}
	
	fmt.Printf("Processed results: %v\n", results)
}