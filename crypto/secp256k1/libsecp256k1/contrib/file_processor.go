package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu      sync.Mutex
	errors  []error
	results []string
}

func (fp *FileProcessor) ProcessFile(path string, info fs.FileInfo, err error) error {
	if err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Errorf("access error %s: %w", path, err))
		fp.mu.Unlock()
		return nil
	}

	if info.IsDir() {
		return nil
	}

	ext := filepath.Ext(path)
	if ext != ".txt" && ext != ".log" {
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fp.mu.Lock()
		fp.errors = append(fp.errors, fmt.Errorf("read error %s: %w", path, err))
		fp.mu.Unlock()
		return nil
	}

	if len(content) > 0 {
		fp.mu.Lock()
		fp.results = append(fp.results, fmt.Sprintf("Processed %s (%d bytes)", path, len(content)))
		fp.mu.Unlock()
	}

	return nil
}

func (fp *FileProcessor) GetResults() ([]string, []error) {
	return fp.results, fp.errors
}

func ProcessDirectory(root string) (*FileProcessor, error) {
	if _, err := os.Stat(root); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("directory does not exist: %s", root)
	}

	fp := &FileProcessor{}
	err := filepath.Walk(root, fp.ProcessFile)
	if err != nil {
		return nil, fmt.Errorf("walk error: %w", err)
	}

	return fp, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	processor, err := ProcessDirectory(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	results, errors := processor.GetResults()

	fmt.Printf("Processed %d files\n", len(results))
	for _, result := range results {
		fmt.Println(result)
	}

	if len(errors) > 0 {
		fmt.Printf("\nEncountered %d errors:\n", len(errors))
		for _, e := range errors {
			fmt.Printf("  - %v\n", e)
		}
	}
}