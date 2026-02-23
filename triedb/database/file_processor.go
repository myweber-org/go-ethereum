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

func (fp *FileProcessor) ProcessFiles(paths []string, processor func(string) error) error {
	if len(paths) == 0 {
		return errors.New("no files to process")
	}

	fileChan := make(chan string, len(paths))
	resultChan := make(chan error, len(paths))

	for i := 0; i < fp.workers; i++ {
		fp.wg.Add(1)
		go fp.worker(fileChan, resultChan, processor)
	}

	for _, path := range paths {
		fileChan <- path
	}
	close(fileChan)

	fp.wg.Wait()
	close(resultChan)

	for err := range resultChan {
		if err != nil {
			return fmt.Errorf("processing error: %w", err)
		}
	}
	return nil
}

func (fp *FileProcessor) worker(files <-chan string, results chan<- error, processor func(string) error) {
	defer fp.wg.Done()

	batch := make([]string, 0, fp.batchSize)
	for file := range files {
		batch = append(batch, file)

		if len(batch) >= fp.batchSize {
			for _, f := range batch {
				results <- processor(f)
			}
			batch = batch[:0]
		}
	}

	for _, f := range batch {
		results <- processor(f)
	}
}

func (fp *FileProcessor) CountLines(path string) (int, error) {
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

func (fp *FileProcessor) CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func main() {
	processor := NewFileProcessor(4, 5)

	sampleFiles := []string{
		"data/file1.txt",
		"data/file2.txt",
		"data/file3.txt",
	}

	processFunc := func(path string) error {
		start := time.Now()
		lines, err := processor.CountLines(path)
		if err != nil {
			return err
		}
		duration := time.Since(start)
		fmt.Printf("Processed %s: %d lines in %v\n", path, lines, duration)
		return nil
	}

	if err := processor.ProcessFiles(sampleFiles, processFunc); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}