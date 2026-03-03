package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu          sync.Mutex
	processed   int
	errors      []error
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		errors: make([]error, 0),
	}
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.recordError(fmt.Errorf("failed to open %s: %w", path, err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fp.recordError(fmt.Errorf("scan error in %s: %w", path, err))
		return
	}

	fp.mu.Lock()
	fp.processed++
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
}

func (fp *FileProcessor) recordError(err error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.errors = append(fp.errors, err)
}

func (fp *FileProcessor) Stats() (int, []error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed, fp.errors
}

func ProcessFiles(paths []string, maxWorkers int) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	if maxWorkers < 1 {
		maxWorkers = 1
	}

	processor := NewFileProcessor()
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, maxWorkers)

	for _, path := range paths {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(p string) {
			defer func() { <-semaphore }()
			processor.ProcessFile(p, &wg)
		}(path)
	}

	wg.Wait()

	processed, errs := processor.Stats()
	fmt.Printf("\nProcessing complete. Files: %d, Errors: %d\n", processed, len(errs))

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during processing", len(errs))
	}

	return nil
}

func main() {
	files := []string{
		"data/file1.txt",
		"data/file2.txt",
		"data/file3.txt",
	}

	if err := ProcessFiles(files, 3); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}