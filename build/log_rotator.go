package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func rotateLog(logPath string, maxBackups int) error {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil
	}

	dir := filepath.Dir(logPath)
	baseName := filepath.Base(logPath)
	timestamp := time.Now().Format("20060102_150405")
	archiveName := fmt.Sprintf("%s.%s", baseName, timestamp)
	archivePath := filepath.Join(dir, archiveName)

	err := os.Rename(logPath, archivePath)
	if err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	pattern := fmt.Sprintf("%s.*", baseName)
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return fmt.Errorf("failed to list backup files: %w", err)
	}

	if len(matches) > maxBackups {
		oldest := matches[:len(matches)-maxBackups]
		for _, oldFile := range oldest {
			if err := os.Remove(oldFile); err != nil {
				return fmt.Errorf("failed to remove old backup %s: %w", oldFile, err)
			}
		}
	}

	return nil
}

func main() {
	logFile := "./app.log"
	if err := rotateLog(logFile, 5); err != nil {
		fmt.Printf("Log rotation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Log rotation completed successfully")
}