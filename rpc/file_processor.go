package main

import (
	"bufio"
	"fmt"
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
}

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", path, err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Printf("Error stating %s: %v\n", path, err)
		return
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	stats := FileStats{
		Path:     path,
		Size:     stat.Size(),
		Lines:    lineCount,
		Modified: stat.ModTime(),
	}

	results <- stats
}

func collectResults(results <-chan FileStats, done chan<- bool) {
	totalFiles := 0
	totalSize := int64(0)
	totalLines := 0

	for stats := range results {
		totalFiles++
		totalSize += stats.Size
		totalLines += stats.Lines
		fmt.Printf("Processed: %s (Size: %d bytes, Lines: %d, Modified: %s)\n",
			filepath.Base(stats.Path), stats.Size, stats.Lines, stats.Modified.Format("2006-01-02"))
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	done <- true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	var wg sync.WaitGroup
	results := make(chan FileStats, 100)
	done := make(chan bool)

	go collectResults(results, done)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		wg.Add(1)
		go processFile(path, &wg, results)
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	wg.Wait()
	close(results)
	<-done
}