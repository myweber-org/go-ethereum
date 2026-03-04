package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	Workers   int
	BatchSize int
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		Workers:   workers,
		BatchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, processor func(string) error) []error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(paths))
	fileChan := make(chan string, fp.BatchSize)

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for path := range fileChan {
				fmt.Printf("Worker %d processing: %s\n", workerID, filepath.Base(path))
				if err := processor(path); err != nil {
					errorChan <- fmt.Errorf("file %s: %w", path, err)
				}
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	wg.Wait()
	close(errorChan)

	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	return errors
}

func sampleProcessor(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	fmt.Printf("File %s has %d lines\n", filepath.Base(filePath), lineCount)
	return scanner.Err()
}

func main() {
	files := []string{
		"/tmp/test1.txt",
		"/tmp/test2.txt",
		"/tmp/test3.txt",
	}

	processor := NewFileProcessor(3, 10)
	errors := processor.ProcessFiles(files, sampleProcessor)

	if len(errors) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errors))
		for _, err := range errors {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("All files processed successfully")
	}
}