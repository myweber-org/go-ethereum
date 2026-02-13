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
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, processor func(string) error) []error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(paths))
	pathChan := make(chan string, fp.batchSize)

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathChan {
				if err := processor(path); err != nil {
					errChan <- fmt.Errorf("processing %s: %w", path, err)
				}
			}
		}()
	}

	for _, path := range paths {
		pathChan <- path
	}
	close(pathChan)

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	return errors
}

func validateFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return errors.New("is a directory")
	}

	if info.Size() == 0 {
		return errors.New("file is empty")
	}

	return nil
}

func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount, scanner.Err()
}

func processFile(path string) error {
	if err := validateFile(path); err != nil {
		return err
	}

	lines, err := countLines(path)
	if err != nil {
		return err
	}

	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lines)
	return nil
}

func main() {
	processor := NewFileProcessor(4, 10)

	files := []string{
		"data/file1.txt",
		"data/file2.txt",
		"data/file3.txt",
	}

	start := time.Now()
	errors := processor.ProcessFiles(files, processFile)
	elapsed := time.Since(start)

	fmt.Printf("Processing completed in %v\n", elapsed)
	if len(errors) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
	}
}