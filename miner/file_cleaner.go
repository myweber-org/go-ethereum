package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/myapp"
	maxAgeDays   = 7
	checkPattern = "*.tmp"
)

func main() {
	err := cleanOldFiles(tempDir, checkPattern, maxAgeDays)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dir, pattern string, maxAgeDays int) error {
	cutoffTime := time.Now().AddDate(0, 0, -maxAgeDays)

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return err
		}

		if matched && info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
		}

		return nil
	})
}