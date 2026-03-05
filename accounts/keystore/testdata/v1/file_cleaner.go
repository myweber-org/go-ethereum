package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	cutoff := time.Now().AddDate(0, 0, -7)
	var removedCount int

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err == nil {
				removedCount++
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	fmt.Printf("Cleaned %d temporary files older than 7 days\n", removedCount)
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	maxAgeHours = 24 * 7 // 1 week
)

func main() {
	dir := os.TempDir()
	fmt.Printf("Cleaning temporary files older than %d hours in: %s\n", maxAgeHours, dir)

	err := cleanOldFiles(dir)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Cleanup completed successfully.")
}

func cleanOldFiles(dirPath string) error {
	cutoffTime := time.Now().Add(-time.Hour * maxAgeHours)

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
		}

		return nil
	})
}