
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileResult struct {
	Path     string
	Size     int64
	Lines    int
	Error    error
}

func processFile(path string, results chan<- FileResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	result := FileResult{Path: path}
	
	file, err := os.Open(path)
	if err != nil {
		result.Error = err
		results <- result
		return
	}
	defer file.Close()
	
	info, err := file.Stat()
	if err != nil {
		result.Error = err
		results <- result
		return
	}
	
	result.Size = info.Size()
	
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	
	if err := scanner.Err(); err != nil {
		result.Error = err
		results <- result
		return
	}
	
	result.Lines = lineCount
	results <- result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}
	
	root := os.Args[1]
	results := make(chan FileResult)
	var wg sync.WaitGroup
	
	startTime := time.Now()
	
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.Mode().IsRegular() {
			return nil
		}
		
		wg.Add(1)
		go processFile(path, results, &wg)
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	totalFiles := 0
	totalSize := int64(0)
	totalLines := 0
	
	for result := range results {
		totalFiles++
		if result.Error != nil {
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Error)
			continue
		}
		
		totalSize += result.Size
		totalLines += result.Lines
		fmt.Printf("Processed: %s (Size: %d bytes, Lines: %d)\n", 
			result.Path, result.Size, result.Lines)
	}
	
	duration := time.Since(startTime)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Processing time: %v\n", duration)
}