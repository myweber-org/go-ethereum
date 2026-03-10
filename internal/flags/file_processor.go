package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func readConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func writeConfig(filename string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func main() {
	config := &Config{
		Host: "localhost",
		Port: 8080,
	}

	err := writeConfig("config.json", config)
	if err != nil {
		fmt.Printf("Error writing config: %v\n", err)
		os.Exit(1)
	}

	loadedConfig, err := readConfig("config.json")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded config: %+v\n", loadedConfig)
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadJSONFile(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func WriteJSONFile(filename string, v interface{}, indent bool) error {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func ValidateJSONStructure(data []byte, expectedFields []string) bool {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}

	for _, field := range expectedFields {
		if _, exists := m[field]; !exists {
			return false
		}
	}
	return true
}
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileResult struct {
	Path    string
	Size    int64
	Lines   int
	Error   error
}

func processFile(path string, results chan<- FileResult, wg *sync.WaitGroup) {
	defer wg.Done()

	result := FileResult{Path: path}
	
	file, err := os.Open(path)
	if err != nil {
		result.Error = err
		results <- result
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		result.Error = err
		results <- result
		return
	}
	result.Size = info.Size()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	
	if err := scanner.Err(); err != nil {
		result.Error = err
		results <- result
		return
	}
	
	result.Lines = lineCount
	results <- result
}

func collectResults(results <-chan FileResult, totalFiles *int, totalSize *int64, totalLines *int) {
	for result := range results {
		if result.Error != nil {
			fmt.Printf("Error processing %s: %v\n", result.Path, result.Error)
			continue
		}
		
		*totalFiles++
		*totalSize += result.Size
		*totalLines += result.Lines
		
		fmt.Printf("Processed: %s (Size: %d bytes, Lines: %d)\n", 
			filepath.Base(result.Path), result.Size, result.Lines)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_processor <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	results := make(chan FileResult, 100)
	var wg sync.WaitGroup

	totalFiles := 0
	var totalSize int64 = 0
	totalLines := 0

	go collectResults(results, &totalFiles, &totalSize, &totalLines)

	startTime := time.Now()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.Mode().IsRegular() {
			return nil
		}

		wg.Add(1)
		go processFile(path, results, &wg)
		
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	wg.Wait()
	close(results)

	duration := time.Since(startTime)
	
	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total files processed: %d\n", totalFiles)
	fmt.Printf("Total size: %d bytes\n", totalSize)
	fmt.Printf("Total lines: %d\n", totalLines)
	fmt.Printf("Processing time: %v\n", duration)
}