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
	workers    int
	queue      chan string
	wg         sync.WaitGroup
	results    map[string]int
	resultLock sync.RWMutex
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		workers: workers,
		queue:   make(chan string, 100),
		results: make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}

	return lineCount, nil
}

func (fp *FileProcessor) worker(id int) {
	defer fp.wg.Done()
	for path := range fp.queue {
		count, err := fp.ProcessFile(path)
		if err != nil {
			fmt.Printf("Worker %d: error processing %s: %v\n", id, path, err)
			continue
		}

		fp.resultLock.Lock()
		fp.results[path] = count
		fp.resultLock.Unlock()

		fmt.Printf("Worker %d: processed %s (%d lines)\n", id, path, count)
		time.Sleep(50 * time.Millisecond)
	}
}

func (fp *FileProcessor) AddFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	fp.queue <- path
	return nil
}

func (fp *FileProcessor) Start() {
	for i := 0; i < fp.workers; i++ {
		fp.wg.Add(1)
		go fp.worker(i + 1)
	}
}

func (fp *FileProcessor) Wait() {
	close(fp.queue)
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() map[string]int {
	fp.resultLock.RLock()
	defer fp.resultLock.RUnlock()

	copied := make(map[string]int)
	for k, v := range fp.results {
		copied[k] = v
	}
	return copied
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	processor := NewFileProcessor(4)
	processor.Start()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			if err := processor.AddFile(path); err != nil {
				fmt.Printf("Error adding file: %v\n", err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	processor.Wait()

	results := processor.GetResults()
	fmt.Printf("\nProcessing complete. Files processed: %d\n", len(results))
	for file, count := range results {
		fmt.Printf("%s: %d lines\n", filepath.Base(file), count)
	}
}