
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	daysToKeep = 7
	tempDir    = "/tmp"
)

func main() {
	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)
	
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if info.IsDir() {
			return nil
		}
		
		if info.ModTime().Before(cutoffTime) {
			fmt.Printf("Removing old file: %s\n", path)
			os.Remove(path)
		}
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("Error cleaning files: %v\n", err)
	}
}