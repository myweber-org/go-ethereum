package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	cutoffTime := time.Now().AddDate(0, 0, -7)
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.ModTime().Before(cutoffTime) {
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("Failed to remove %s: %v\n", path, err)
			} else {
				fmt.Printf("Removed: %s\n", path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}