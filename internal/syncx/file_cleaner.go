package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/myapp"
	maxAgeHours  = 168
)

func main() {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	now := time.Now()
	removedCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileAge := now.Sub(file.ModTime())
		if fileAge.Hours() > maxAgeHours {
			filePath := filepath.Join(tempDir, file.Name())
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", file.Name(), err)
			} else {
				removedCount++
				fmt.Printf("Removed old file: %s\n", file.Name())
			}
		}
	}

	fmt.Printf("Cleanup completed. Removed %d files.\n", removedCount)
}