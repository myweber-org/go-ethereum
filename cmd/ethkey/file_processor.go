
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
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				fp.processFile(p)
			}(path)
		}
		return nil
	})
	wg.Wait()
	return err
}

func (fp *FileProcessor) processFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
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
	err := processor.ProcessDirectory(os.Args[1])
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}

	results := processor.GetResults()
	for _, result := range results {
		fmt.Println(result)
	}
}