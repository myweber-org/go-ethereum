package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileProcessor struct {
	InputDir  string
	OutputDir string
	Workers   int
}

func NewFileProcessor(inputDir, outputDir string, workers int) *FileProcessor {
	return &FileProcessor{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Workers:   workers,
	}
}

func (fp *FileProcessor) ProcessFiles() error {
	files, err := filepath.Glob(filepath.Join(fp.InputDir, "*.txt"))
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	var wg sync.WaitGroup
	fileChan := make(chan string, len(files))

	for i := 0; i < fp.Workers; i++ {
		wg.Add(1)
		go fp.worker(fileChan, &wg)
	}

	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	wg.Wait()
	return nil
}

func (fp *FileProcessor) worker(files <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range files {
		if err := fp.processSingleFile(file); err != nil {
			fmt.Printf("error processing %s: %v\n", file, err)
		}
	}
}

func (fp *FileProcessor) processSingleFile(inputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	outputPath := filepath.Join(fp.OutputDir, filepath.Base(inputPath))
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		processed := transformLine(line)
		if _, err := writer.WriteString(processed + "\n"); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func transformLine(line string) string {
	return fmt.Sprintf("PROCESSED: %s", line)
}

func main() {
	processor := NewFileProcessor("./input", "./output", 4)
	if err := processor.ProcessFiles(); err != nil {
		fmt.Printf("processing failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("file processing completed successfully")
}