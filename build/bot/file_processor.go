
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	Workers    int
	BatchSize  int
	ResultChan chan ProcessResult
	ErrorChan  chan error
}

type ProcessResult struct {
	Filename string
	Lines    int
	Duration time.Duration
}

func NewFileProcessor(workers, batchSize int) *FileProcessor {
	return &FileProcessor{
		Workers:    workers,
		BatchSize:  batchSize,
		ResultChan: make(chan ProcessResult, 100),
		ErrorChan:  make(chan error, 100),
	}
}

func (fp *FileProcessor) ProcessFiles(paths []string) {
	var wg sync.WaitGroup
	fileBatches := fp.createBatches(paths)

	for i, batch := range fileBatches {
		wg.Add(1)
		go func(batchIndex int, batchFiles []string) {
			defer wg.Done()
			fp.processBatch(batchIndex, batchFiles)
		}(i, batch)
	}

	go func() {
		wg.Wait()
		close(fp.ResultChan)
		close(fp.ErrorChan)
	}()
}

func (fp *FileProcessor) createBatches(paths []string) [][]string {
	var batches [][]string
	for i := 0; i < len(paths); i += fp.BatchSize {
		end := i + fp.BatchSize
		if end > len(paths) {
			end = len(paths)
		}
		batches = append(batches, paths[i:end])
	}
	return batches
}

func (fp *FileProcessor) processBatch(batchIndex int, files []string) {
	semaphore := make(chan struct{}, fp.Workers)
	var batchWg sync.WaitGroup

	for _, filepath := range files {
		batchWg.Add(1)
		semaphore <- struct{}{}

		go func(fpath string) {
			defer func() {
				<-semaphore
				batchWg.Done()
			}()

			result, err := fp.processSingleFile(fpath)
			if err != nil {
				fp.ErrorChan <- fmt.Errorf("file %s: %w", fpath, err)
				return
			}
			fp.ResultChan <- result
		}(filepath)
	}
	batchWg.Wait()
}

func (fp *FileProcessor) processSingleFile(filepath string) (ProcessResult, error) {
	start := time.Now()

	file, err := os.Open(filepath)
	if err != nil {
		return ProcessResult{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return ProcessResult{}, err
	}

	duration := time.Since(start)
	return ProcessResult{
		Filename: filepath,
		Lines:    lineCount,
		Duration: duration,
	}, nil
}

func validateFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("is a directory")
	}
	if info.Size() == 0 {
		return errors.New("file is empty")
	}
	return nil
}

func collectResults(resultChan <-chan ProcessResult, errorChan <-chan error) {
	for {
		select {
		case result, ok := <-resultChan:
			if !ok {
				resultChan = nil
			} else {
				fmt.Printf("Processed: %s, Lines: %d, Time: %v\n",
					filepath.Base(result.Filename),
					result.Lines,
					result.Duration)
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else {
				fmt.Printf("Error: %v\n", err)
			}
		}

		if resultChan == nil && errorChan == nil {
			break
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <file1> <file2> ...")
		os.Exit(1)
	}

	filePaths := os.Args[1:]
	var validPaths []string

	for _, path := range filePaths {
		if err := validateFile(path); err != nil {
			fmt.Printf("Skipping %s: %v\n", path, err)
			continue
		}
		validPaths = append(validPaths, path)
	}

	if len(validPaths) == 0 {
		fmt.Println("No valid files to process")
		return
	}

	processor := NewFileProcessor(4, 10)
	go processor.ProcessFiles(validPaths)

	collectResults(processor.ResultChan, processor.ErrorChan)
	fmt.Println("File processing completed")
}