package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type FileProcessor struct {
	inputPath  string
	outputPath string
	mu         sync.Mutex
}

func NewFileProcessor(input, output string) *FileProcessor {
	return &FileProcessor{
		inputPath:  input,
		outputPath: output,
	}
}

func (fp *FileProcessor) ProcessLines() error {
	file, err := os.Open(fp.inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	outputFile, err := os.Create(fp.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	scanner := bufio.NewScanner(file)
	var wg sync.WaitGroup
	lineChan := make(chan string, 100)
	results := make(chan string, 100)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go fp.worker(i, lineChan, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for scanner.Scan() {
			lineChan <- scanner.Text()
		}
		close(lineChan)
	}()

	for result := range results {
		_, err := writer.WriteString(result + "\n")
		if err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

func (fp *FileProcessor) worker(id int, lines <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range lines {
		processed := fmt.Sprintf("Worker %d: %s", id, line)
		fp.mu.Lock()
		results <- processed
		fp.mu.Unlock()
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: file_processor <input_file> <output_file>")
		os.Exit(1)
	}

	processor := NewFileProcessor(os.Args[1], os.Args[2])
	if err := processor.ProcessLines(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File processing completed successfully")
}