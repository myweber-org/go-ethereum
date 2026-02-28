package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type FileProcessor struct {
	mu       sync.Mutex
	results  map[string]int
	wg       sync.WaitGroup
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		results: make(map[string]int),
	}
}

func (fp *FileProcessor) ProcessFile(filename string) error {
	file, err := os.Open(filename)
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
	fp.results[filename] = lineCount
	fp.mu.Unlock()

	return scanner.Err()
}

func (fp *FileProcessor) ConcurrentProcess(files []string) {
	for _, file := range files {
		fp.wg.Add(1)
		go func(f string) {
			defer fp.wg.Done()
			if err := fp.ProcessFile(f); err != nil {
				fmt.Printf("Error processing %s: %v\n", f, err)
			}
		}(file)
	}
	fp.wg.Wait()
}

func (fp *FileProcessor) GetResults() map[string]int {
	fp.mu.Lock()
	defer fp.mu.Unlock()
	return fp.results
}

func main() {
	files := []string{"data1.txt", "data2.txt", "data3.txt"}
	processor := NewFileProcessor()
	processor.ConcurrentProcess(files)
	results := processor.GetResults()
	for file, count := range results {
		fmt.Printf("%s: %d lines\n", file, count)
	}
}
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

type RecordProcessor struct {
	mu       sync.Mutex
	processed int
}

func (rp *RecordProcessor) Process(record []string) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	
	rp.processed++
	fmt.Printf("Processing record %d: %v\n", rp.processed, record)
}

func main() {
	file, err := os.Open("data.csv")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	processor := &RecordProcessor{}
	var wg sync.WaitGroup

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading CSV: %v\n", err)
			continue
		}

		wg.Add(1)
		go func(r []string) {
			defer wg.Done()
			processor.Process(r)
		}(record)
	}

	wg.Wait()
	fmt.Printf("Total records processed: %d\n", processor.processed)
}