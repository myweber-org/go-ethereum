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
	mu        sync.Mutex
	wg        sync.WaitGroup
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	if workers < 1 {
		workers = 3
	}
	if batchSize < 1 {
		batchSize = 100
	}
	return &FileProcessor{
		workers:   workers,
		batchSize: batchSize,
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string, handler func(string) error) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	fileChan := make(chan string, fp.batchSize)
	errChan := make(chan error, fp.workers)
	done := make(chan struct{})

	for i := 0; i < fp.workers; i++ {
		fp.wg.Add(1)
		go fp.worker(fileChan, errChan, handler)
	}

	go func() {
		for _, path := range paths {
			absPath, err := filepath.Abs(path)
			if err != nil {
				errChan <- fmt.Errorf("invalid path %s: %w", path, err)
				continue
			}
			fileChan <- absPath
		}
		close(fileChan)
	}()

	go func() {
		fp.wg.Wait()
		close(done)
	}()

	var processErr error
	select {
	case err := <-errChan:
		if processErr == nil {
			processErr = err
		}
	case <-done:
	}

	return processErr
}

func (fp *FileProcessor) worker(files <-chan string, errChan chan<- error, handler func(string) error) {
	defer fp.wg.Done()
	for file := range files {
		if err := handler(file); err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}
}

func countLines(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", path, err)
	}

	fp.mu.Lock()
	fmt.Printf("File %s has %d lines\n", filepath.Base(path), lineCount)
	fp.mu.Unlock()
	return nil
}

func main() {
	processor := NewFileProcessor(4, 50)

	sampleFiles := []string{
		"/tmp/test1.txt",
		"/tmp/test2.log",
		"/tmp/data.json",
	}

	start := time.Now()
	err := processor.ProcessFiles(sampleFiles, countLines)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Processing error: %v\n", err)
	}

	fmt.Printf("Processed in %v\n", elapsed)
}