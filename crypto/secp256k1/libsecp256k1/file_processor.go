package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	workers   int
	batchSize int
	mu        sync.Mutex
	results   map[string]int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
		results:   make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return errors.New("directory does not exist")
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, fp.batchSize)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(&wg, fileChan)
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileChan <- path
		}
		return nil
	})

	close(fileChan)
	wg.Wait()

	if err != nil {
		return fmt.Errorf("walk error: %v", err)
	}

	return nil
}

func (fp *FileProcessor) worker(wg *sync.WaitGroup, fileChan <-chan string) {
	defer wg.Done()

	for filePath := range fileChan {
		count, err := fp.countLines(filePath)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
			continue
		}

		fp.mu.Lock()
		fp.results[filePath] = count
		fp.mu.Unlock()
	}
}

func (fp *FileProcessor) countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

func (fp *FileProcessor) PrintSummary() {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	fmt.Println("Processing Summary:")
	fmt.Println("===================")
	totalFiles := len(fp.results)
	totalLines := 0

	for file, count := range fp.results {
		fmt.Printf("%s: %d lines\n", filepath.Base(file), count)
		totalLines += count
	}

	fmt.Printf("\nTotal files processed: %d\n", totalFiles)
	fmt.Printf("Total lines counted: %d\n", totalLines)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_processor.go <directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor(4, 100)

	startTime := time.Now()
	err := processor.ProcessDirectory(dirPath)
	elapsed := time.Since(startTime)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	processor.PrintSummary()
	fmt.Printf("\nProcessing completed in %v\n", elapsed)
}