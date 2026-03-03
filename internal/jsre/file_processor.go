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
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Printf("Error stating file %s: %v\n", path, err)
		return
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return
	}

	stats := FileStats{
		Path:     path,
		Size:     info.Size(),
		Lines:    lineCount,
		Modified: info.ModTime(),
	}

	results <- stats
}

func collectFiles(dir string, patterns []string) ([]string, error) {
	var files []string

	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}

	return files, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory> [patterns...]")
		fmt.Println("Default patterns: *.txt, *.go, *.md")
		os.Exit(1)
	}

	dir := os.Args[1]
	patterns := []string{"*.txt", "*.go", "*.md"}

	if len(os.Args) > 2 {
		patterns = os.Args[2:]
	}

	files, err := collectFiles(dir, patterns)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found matching patterns")
		return
	}

	var wg sync.WaitGroup
	results := make(chan FileStats, len(files))

	for _, file := range files {
		wg.Add(1)
		go processFile(file, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalSize := int64(0)
	totalLines := 0
	fileCount := 0

	fmt.Println("File Processing Results:")
	fmt.Println("========================")

	for stats := range results {
		fmt.Printf("File: %s\n", stats.Path)
		fmt.Printf("  Size: %d bytes\n", stats.Size)
		fmt.Printf("  Lines: %d\n", stats.Lines)
		fmt.Printf("  Modified: %s\n", stats.Modified.Format("2006-01-02 15:04:05"))
		fmt.Println()

		totalSize += stats.Size
		totalLines += stats.Lines
		fileCount++
	}

	fmt.Println("Summary:")
	fmt.Printf("  Total files processed: %d\n", fileCount)
	fmt.Printf("  Total size: %d bytes\n", totalSize)
	fmt.Printf("  Total lines: %d\n", totalLines)
}