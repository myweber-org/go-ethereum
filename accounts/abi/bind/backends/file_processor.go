
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	workers   int
	batchSize int
	mu        sync.RWMutex
	stats     ProcessingStats
}

type ProcessingStats struct {
	FilesProcessed int
	TotalBytes     int64
	Errors         int
	StartTime      time.Time
	EndTime        time.Time
}

type FileTask struct {
	Path    string
	Content []byte
	Err     error
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	if workers < 1 {
		workers = 4
	}
	if batchSize < 1 {
		batchSize = 100
	}

	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		stats:     ProcessingStats{},
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	fp.mu.Lock()
	fp.stats = ProcessingStats{StartTime: time.Now()}
	fp.mu.Unlock()

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	taskChan := make(chan FileTask, fp.batchSize)
	resultChan := make(chan FileTask, fp.batchSize)
	var wg sync.WaitGroup

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(taskChan, resultChan, &wg)
	}

	go fp.collectResults(resultChan)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullPath := filepath.Join(dirPath, file.Name())
		taskChan <- FileTask{Path: fullPath}
	}

	close(taskChan)
	wg.Wait()
	close(resultChan)

	fp.mu.Lock()
	fp.stats.EndTime = time.Now()
	fp.mu.Unlock()

	return nil
}

func (fp *FileProcessor) worker(tasks <-chan FileTask, results chan<- FileTask, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		content, err := fp.readFile(task.Path)
		task.Content = content
		task.Err = err
		results <- task
	}
}

func (fp *FileProcessor) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var content []byte
	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		if n == 0 {
			break
		}

		content = append(content, buffer[:n]...)
	}

	return content, nil
}

func (fp *FileProcessor) collectResults(results <-chan FileTask) {
	for result := range results {
		fp.mu.Lock()
		fp.stats.FilesProcessed++

		if result.Err != nil {
			fp.stats.Errors++
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Err)
		} else {
			fp.stats.TotalBytes += int64(len(result.Content))
		}
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) GetStats() ProcessingStats {
	fp.mu.RLock()
	defer fp.mu.RUnlock()
	return fp.stats
}

func (fp *FileProcessor) PrintStats() {
	stats := fp.GetStats()
	duration := stats.EndTime.Sub(stats.StartTime)

	fmt.Println("\n=== Processing Statistics ===")
	fmt.Printf("Files processed: %d\n", stats.FilesProcessed)
	fmt.Printf("Total bytes: %d\n", stats.TotalBytes)
	fmt.Printf("Errors: %d\n", stats.Errors)
	fmt.Printf("Duration: %v\n", duration.Round(time.Millisecond))

	if stats.FilesProcessed > 0 && duration > 0 {
		throughput := float64(stats.TotalBytes) / duration.Seconds()
		fmt.Printf("Throughput: %.2f bytes/sec\n", throughput)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory_path>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor(4, 100)

	fmt.Printf("Processing directory: %s\n", dirPath)
	err := processor.ProcessDirectory(dirPath)
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	processor.PrintStats()
}