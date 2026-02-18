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
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	if workers < 1 {
		workers = 1
	}
	if batchSize < 1 {
		batchSize = 10
	}
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, fp.batchSize)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go fp.worker(i, fileChan, &wg)
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(dirPath, file.Name())
			fileChan <- fullPath
		}
	}

	close(fileChan)
	wg.Wait()

	return nil
}

func (fp *FileProcessor) worker(id int, files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for filePath := range files {
		if err := fp.processFile(filePath); err != nil {
			fmt.Printf("Worker %d: error processing %s: %v\n", id, filePath, err)
		} else {
			fmt.Printf("Worker %d: successfully processed %s\n", id, filePath)
		}
	}
}

func (fp *FileProcessor) processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		time.Sleep(1 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	if lineCount == 0 {
		return errors.New("empty file")
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	processor := NewFileProcessor(4, 20)

	start := time.Now()
	if err := processor.ProcessDirectory(dirPath); err != nil {
		fmt.Printf("Processing failed: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(start)
	fmt.Printf("Processing completed in %v\n", duration)
}