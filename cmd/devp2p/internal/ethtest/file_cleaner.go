package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	tempFilePrefix = "temp_"
	maxAgeHours    = 24
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	err := cleanTempFiles(dir)
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Cleanup completed successfully")
}

func cleanTempFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !isTempFile(info.Name()) {
			return nil
		}

		if isFileOld(info.ModTime()) {
			return os.Remove(path)
		}

		return nil
	})
}

func isTempFile(filename string) bool {
	return len(filename) > len(tempFilePrefix) && filename[:len(tempFilePrefix)] == tempFilePrefix
}

func isFileOld(modTime time.Time) bool {
	return time.Since(modTime) > maxAgeHours*time.Hour
}