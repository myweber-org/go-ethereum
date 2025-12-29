package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileProcessor struct {
	workers    int
	bufferSize int
	verbose    bool
}

func NewFileProcessor(workers, bufferSize int, verbose bool) *FileProcessor {
	if workers < 1 {
		workers = 1
	}
	if bufferSize < 1024 {
		bufferSize = 4096
	}
	return &FileProcessor{
		workers:    workers,
		bufferSize: bufferSize,
		verbose:    verbose,
	}
}

func (fp *FileProcessor) ProcessFile(path string, processor func([]byte) ([]byte, error)) error {
	if !fileExists(path) {
		return errors.New("file does not exist")
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	outputPath := path + ".processed"
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	reader := bufio.NewReaderSize(file, fp.bufferSize)
	writer := bufio.NewWriterSize(outFile, fp.bufferSize)
	defer writer.Flush()

	var wg sync.WaitGroup
	chunkChan := make(chan []byte, fp.workers*2)
	errChan := make(chan error, 1)
	done := make(chan bool)

	go func() {
		for {
			chunk := make([]byte, fp.bufferSize)
			n, err := reader.Read(chunk)
			if n > 0 {
				chunkChan <- chunk[:n]
			}
			if err != nil {
				if err != io.EOF {
					errChan <- err
				}
				close(chunkChan)
				break
			}
		}
	}()

	for i := 0; i < fp.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for chunk := range chunkChan {
				processed, err := processor(chunk)
				if err != nil {
					errChan <- fmt.Errorf("worker %d: %w", workerID, err)
					return
				}
				if _, err := writer.Write(processed); err != nil {
					errChan <- fmt.Errorf("worker %d write failed: %w", workerID, err)
					return
				}
				if fp.verbose {
					fmt.Printf("Worker %d processed %d bytes\n", workerID, len(chunk))
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:
		return err
	case <-done:
		if fp.verbose {
			fmt.Println("File processing completed successfully")
		}
		return nil
	case <-time.After(30 * time.Second):
		return errors.New("processing timeout")
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func main() {
	processor := NewFileProcessor(4, 8192, true)

	sampleProcessor := func(data []byte) ([]byte, error) {
		for i := range data {
			data[i] = data[i] ^ 0xFF
		}
		return data, nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	testFile := filepath.Join(currentDir, "test_input.txt")
	if !fileExists(testFile) {
		content := []byte("Sample content for processing test\n")
		if err := os.WriteFile(testFile, content, 0644); err != nil {
			fmt.Printf("Error creating test file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created test file:", testFile)
	}

	if err := processor.ProcessFile(testFile, sampleProcessor); err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing completed. Output file:", testFile+".processed")
}