package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_temp"
	maxAgeDays   = 7
	filePattern  = "*.tmp"
)

func main() {
	if err := cleanOldFiles(); err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles() error {
	files, err := filepath.Glob(filepath.Join(tempDir, filePattern))
	if err != nil {
		return err
	}

	cutoffTime := time.Now().AddDate(0, 0, -maxAgeDays)
	removedCount := 0

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffTime) {
			if err := os.Remove(file); err == nil {
				removedCount++
			}
		}
	}

	fmt.Printf("Removed %d temporary files older than %d days\n", removedCount, maxAgeDays)
	return nil
}