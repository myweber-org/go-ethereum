package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	dir := "./tmp"
	cleanupAge := 24 * time.Hour

	err := cleanOldFiles(dir, cleanupAge)
	if err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cleanup completed successfully")
}

func cleanOldFiles(dir string, maxAge time.Duration) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if time.Since(info.ModTime()) > maxAge {
			fmt.Printf("Removing old file: %s\n", path)
			return os.Remove(path)
		}
		return nil
	})
}