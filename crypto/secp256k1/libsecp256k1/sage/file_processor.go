
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

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		errors: make([]string, 0),
	}
}

func (fp *FileProcessor) ProcessFile(path string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Sprintf("failed to open %s: %v", path, err))
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
		fp.errors = append(fp.errors, fmt.Sprintf("error scanning %s: %v", path, err))
		fp.mu.Unlock()
		return
	}

	fp.mu.Lock()
	fp.processed++
	fp.mu.Unlock()

	fmt.Printf("Processed %s: %d lines\n", filepath.Base(path), lineCount)
}

func (fp *FileProcessor) Stats() {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	fmt.Printf("\nProcessing complete:\n")
	fmt.Printf("Files processed: %d\n", fp.processed)
	fmt.Printf("Errors encountered: %d\n", len(fp.errors))
	
	if len(fp.errors) > 0 {
		fmt.Println("\nError details:")
		for _, err := range fp.errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> <file2> ...")
		os.Exit(1)
	}

	processor := NewFileProcessor()
	var wg sync.WaitGroup

	for _, filePath := range os.Args[1:] {
		wg.Add(1)
		go processor.ProcessFile(filePath, &wg)
	}

	wg.Wait()
	processor.Stats()
}