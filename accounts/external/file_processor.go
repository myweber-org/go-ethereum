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
		fp.recordError(fmt.Errorf("error scanning %s: %w", path, err))
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

func ProcessFiles(paths []string) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	processor := NewFileProcessor()
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go processor.ProcessFile(path, &wg)
	}

	wg.Wait()

	processed, errs := processor.Stats()
	fmt.Printf("\nProcessing complete: %d files processed\n", processed)

	if len(errs) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("  - %v\n", err)
		}
		return errors.New("some files failed to process")
	}

	return nil
}

func main() {
	files := []string{"data1.txt", "data2.txt", "data3.txt"}
	if err := ProcessFiles(files); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}