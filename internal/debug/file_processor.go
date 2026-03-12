package main

import (
	"fmt"
	"sync"
	"time"
)

type FileProcessor struct {
	workerCount int
	batchSize   int
}

func NewFileProcessor(workers, batch int) *FileProcessor {
	return &FileProcessor{
		workerCount: workers,
		batchSize:   batch,
	}
}

func (fp *FileProcessor) ProcessBatch(fileIDs []string) []string {
	var wg sync.WaitGroup
	results := make([]string, len(fileIDs))
	chunkSize := (len(fileIDs) + fp.workerCount - 1) / fp.workerCount

	for i := 0; i < fp.workerCount; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(fileIDs) {
			end = len(fileIDs)
		}

		if start >= end {
			break
		}

		wg.Add(1)
		go func(workerID int, batch []string) {
			defer wg.Done()
			for idx, fileID := range batch {
				results[start+idx] = fp.processFile(fileID)
			}
		}(i, fileIDs[start:end])
	}

	wg.Wait()
	return results
}

func (fp *FileProcessor) processFile(fileID string) string {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("processed_%s", fileID)
}

func main() {
	processor := NewFileProcessor(4, 25)

	fileIDs := make([]string, 100)
	for i := range fileIDs {
		fileIDs[i] = fmt.Sprintf("file_%03d", i+1)
	}

	start := time.Now()
	results := processor.ProcessBatch(fileIDs)
	elapsed := time.Since(start)

	fmt.Printf("Processed %d files in %v\n", len(results), elapsed)
	for i := 0; i < 3; i++ {
		fmt.Printf("Sample result: %s\n", results[i])
	}
}