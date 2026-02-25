package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu          sync.Mutex
	processed   int
	errors      []string
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("Failed to open %s: %v", path, err))
		fp.mu.Unlock()
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("Error scanning %s: %v", path, err))
		fp.mu.Unlock()
		return
	}

	fp.mu.Lock()
	fp.processed++
	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
	fp.mu.Unlock()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> [file2] ...")
		return
	}

	processor := &FileProcessor{}
	var wg sync.WaitGroup

	for _, filePath := range os.Args[1:] {
		wg.Add(1)
		go processor.ProcessFile(filePath, &wg)
	}

	wg.Wait()

	fmt.Printf("\nSummary: Processed %d files successfully\n", processor.processed)
	if len(processor.errors) > 0 {
		fmt.Printf("Encountered %d errors:\n", len(processor.errors))
		for _, err := range processor.errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}