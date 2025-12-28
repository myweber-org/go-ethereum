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

type FileStats struct {
	Path     string
	Size     int64
	Lines    int
	Modified time.Time
	Error    error
}

type FileProcessor struct {
	workers   int
	results   chan FileStats
	wg        sync.WaitGroup
	mu        sync.Mutex
	processed int
}

func NewFileProcessor(workers int) *FileProcessor {
	return &FileProcessor{
		workers: workers,
		results: make(chan FileStats, 100),
	}
}

func (fp *FileProcessor) ProcessFile(path string) {
	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()

		stats := FileStats{Path: path}
		info, err := os.Stat(path)
		if err != nil {
			stats.Error = err
			fp.results <- stats
			return
		}

		stats.Size = info.Size()
		stats.Modified = info.ModTime()

		file, err := os.Open(path)
		if err != nil {
			stats.Error = err
			fp.results <- stats
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineCount := 0
		for scanner.Scan() {
			lineCount++
		}

		if err := scanner.Err(); err != nil {
			stats.Error = err
		} else {
			stats.Lines = lineCount
		}

		fp.results <- stats
	}()
}

func (fp *FileProcessor) WalkDirectory(root string) error {
	fileQueue := make(chan string, 1000)

	for i := 0; i < fp.workers; i++ {
		go func() {
			for path := range fileQueue {
				fp.ProcessFile(path)
			}
		}()
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		fileQueue <- path
		return nil
	})

	close(fileQueue)
	return err
}

func (fp *FileProcessor) CollectResults() []FileStats {
	go func() {
		fp.wg.Wait()
		close(fp.results)
	}()

	var allStats []FileStats
	for stats := range fp.results {
		fp.mu.Lock()
		fp.processed++
		fp.mu.Unlock()
		allStats = append(allStats, stats)
	}

	return allStats
}

func (fp *FileProcessor) GetProcessedCount() int {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.processed
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	processor := NewFileProcessor(10)

	fmt.Printf("Processing files in: %s\n", dir)
	start := time.Now()

	if err := processor.WalkDirectory(dir); err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	results := processor.CollectResults()
	duration := time.Since(start)

	var totalSize int64
	var totalLines int
	var errorCount int

	for _, stats := range results {
		if stats.Error != nil {
			errorCount++
			fmt.Printf("Error processing %s: %v\n", stats.Path, stats.Error)
			continue
		}

		totalSize += stats.Size
		totalLines += stats.Lines
	}

	fmt.Printf("\nProcessing completed in %v\n", duration)
	fmt.Printf("Files processed: %d\n", processor.GetProcessedCount())
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Errors encountered: %d\n", errorCount)
}