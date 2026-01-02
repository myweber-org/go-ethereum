package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileJob struct {
	Path    string
	Content string
}

func processFile(path string) (FileJob, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileJob{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	return FileJob{Path: path, Content: content}, scanner.Err()
}

func worker(id int, jobs <-chan string, results chan<- FileJob, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range jobs {
		fmt.Printf("Worker %d processing: %s\n", id, filepath.Base(path))
		result, err := processFile(path)
		if err != nil {
			fmt.Printf("Error processing %s: %v\n", path, err)
			continue
		}
		results <- result
	}
}

func main() {
	start := time.Now()
	filePaths := []string{
		"data/file1.txt",
		"data/file2.txt",
		"data/file3.txt",
	}

	numWorkers := 2
	jobs := make(chan string, len(filePaths))
	results := make(chan FileJob, len(filePaths))

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for _, path := range filePaths {
		jobs <- path
	}
	close(jobs)

	wg.Wait()
	close(results)

	var processed []FileJob
	for result := range results {
		processed = append(processed, result)
	}

	fmt.Printf("Processed %d files in %v\n", len(processed), time.Since(start))
	for _, job := range processed {
		fmt.Printf("File: %s, Characters: %d\n", job.Path, len(job.Content))
	}
}