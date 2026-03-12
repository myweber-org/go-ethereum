
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	fileList []string
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		fileList: make([]string, 0),
	}
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				if err := fp.processFile(p); err != nil {
					select {
					case errChan <- err:
					default:
					}
				}
			}(path)
		}
		return nil
	}

	go func() {
		if err := filepath.Walk(dirPath, walkFunc); err != nil {
			errChan <- err
		}
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func (fp *FileProcessor) processFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	fp.mu.Lock()
	fp.fileList = append(fp.fileList, fmt.Sprintf("%s: %d lines", filePath, lineCount))
	fp.mu.Unlock()

	return scanner.Err()
}

func (fp *FileProcessor) GetResults() []string {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.fileList
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	processor := NewFileProcessor()
	if err := processor.ProcessDirectory(os.Args[1]); err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}

	results := processor.GetResults()
	fmt.Printf("Processed %d files:\n", len(results))
	for _, result := range results {
		fmt.Println(result)
	}
}