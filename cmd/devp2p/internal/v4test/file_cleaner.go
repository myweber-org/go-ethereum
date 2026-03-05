package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	tempDir := os.TempDir()
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		fmt.Printf("Error reading temp directory: %v\n", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -7)
	removedCount := 0

	for _, file := range files {
		if file.ModTime().Before(cutoff) {
			filePath := filepath.Join(tempDir, file.Name())
			err := os.RemoveAll(filePath)
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", filePath, err)
			} else {
				removedCount++
			}
		}
	}

	fmt.Printf("Cleaned up %d old temporary files.\n", removedCount)
}package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const (
	tempDir      = "/tmp/app_cache"
	retentionDays = 7
)

func main() {
	err := cleanOldFiles(tempDir, retentionDays)
	if err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dirPath string, days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing: %s (modified: %v)\n", path, info.ModTime())
			return os.Remove(path)
		}

		return nil
	})
}