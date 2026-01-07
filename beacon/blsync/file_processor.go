
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
	ReadTime time.Duration
}

func processFile(path string, wg *sync.WaitGroup, results chan<- FileStats) {
	defer wg.Done()

	start := time.Now()
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", path, err)
		return
	}
	defer file.Close()

	var lineCount int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	fileInfo, _ := file.Stat()
	stats := FileStats{
		Path:     path,
		Size:     fileInfo.Size(),
		Lines:    lineCount,
		ReadTime: time.Since(start),
	}

	results <- stats
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run file_processor.go <directory>")
		return
	}

	root := os.Args[1]
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	results := make(chan FileStats, len(files))

	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg, results)
	}

	wg.Wait()
	close(results)

	totalFiles := 0
	totalLines := 0
	var totalSize int64
	var totalTime time.Duration

	for stats := range results {
		totalFiles++
		totalLines += stats.Lines
		totalSize += stats.Size
		totalTime += stats.ReadTime
		fmt.Printf("Processed: %s | Size: %d bytes | Lines: %d | Time: %v\n",
			stats.Path, stats.Size, stats.Lines, stats.ReadTime)
	}

	fmt.Printf("\nSummary: %d files | %d total lines | %d total bytes | %v total processing time\n",
		totalFiles, totalLines, totalSize, totalTime)
}