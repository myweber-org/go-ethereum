
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	maxFileSize   = 10 * 1024 * 1024 // 10MB
	backupCount   = 5
	logDir        = "./logs"
	currentLog    = "app.log"
	compressOld   = true
)

type LogRotator struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	basePath  string
}

func NewLogRotator() (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	filePath := filepath.Join(logDir, currentLog)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &LogRotator{
		file:     file,
		size:     info.Size(),
		basePath: filePath,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	n, err := lr.file.Write(p)
	if err != nil {
		return n, err
	}

	lr.size += int64(n)
	if lr.size >= maxFileSize {
		if err := lr.rotate(); err != nil {
			log.Printf("rotation failed: %v", err)
		}
	}

	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", timestamp))

	if err := os.Rename(lr.basePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	lr.file = file
	lr.size = 0

	go lr.manageBackups(backupPath)

	return nil
}

func (lr *LogRotator) manageBackups(newBackup string) {
	pattern := filepath.Join(logDir, "app_*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("failed to list backup files: %v", err)
		return
	}

	if len(matches) > backupCount {
		sortByModTime(matches)
		for i := 0; i < len(matches)-backupCount; i++ {
			if compressOld {
				if err := compressFile(matches[i]); err != nil {
					log.Printf("compression failed for %s: %v", matches[i], err)
				}
			} else {
				if err := os.Remove(matches[i]); err != nil {
					log.Printf("failed to remove old log %s: %v", matches[i], err)
				}
			}
		}
	}
}

func sortByModTime(files []string) {
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			infoI, _ := os.Stat(files[i])
			infoJ, _ := os.Stat(files[j])
			if infoI.ModTime().After(infoJ.ModTime()) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
}

func compressFile(path string) error {
	if strings.HasSuffix(path, ".gz") {
		return nil
	}

	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	dstPath := path + ".gz"
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Simple copy for demonstration (replace with actual compression)
	_, err = io.Copy(dst, src)
	if err != nil {
		os.Remove(dstPath)
		return err
	}

	if err := os.Remove(path); err != nil {
		os.Remove(dstPath)
		return err
	}

	return nil
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator()
	if err != nil {
		log.Fatalf("Failed to initialize log rotator: %v", err)
	}
	defer rotator.Close()

	log.SetOutput(rotator)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}