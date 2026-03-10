package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: file_cleaner <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := inputFile + ".cleaned"

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	seen := make(map[string]bool)
	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] {
			seen[line] = true
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

	fmt.Printf("Cleaned file saved as: %s\n", outputFile)
}package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	daysToKeep   = 7
	fileExt      = ".tmp"
)

func main() {
	err := cleanOldFiles(tempDir, daysToKeep, fileExt)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, days int, extension string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -days)
	removedCount := 0

	for _, file := range files {
		if filepath.Ext(file.Name()) != extension {
			continue
		}

		if file.ModTime().Before(cutoffTime) {
			fullPath := filepath.Join(dirPath, file.Name())
			err := os.Remove(fullPath)
			if err != nil {
				fmt.Printf("Warning: Could not remove %s: %v\n", fullPath, err)
				continue
			}
			removedCount++
			fmt.Printf("Removed: %s (modified: %v)\n", file.Name(), file.ModTime())
		}
	}

	fmt.Printf("Total files removed: %d\n", removedCount)
	return nil
}