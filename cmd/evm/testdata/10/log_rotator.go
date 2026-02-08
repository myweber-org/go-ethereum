package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type LogRotator struct {
	mu           sync.Mutex
	filePath     string
	maxSize      int64
	currentSize  int64
	file         *os.File
	rotationCount int
}

func NewLogRotator(filePath string, maxSizeMB int) (*LogRotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		filePath:     filePath,
		maxSize:      maxSize,
		currentSize:  info.Size(),
		file:         file,
		rotationCount: 0,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	lr.file.Close()
	lr.rotationCount++

	backupPath := fmt.Sprintf("%s.%d", lr.filePath, lr.rotationCount)
	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.file = file
	lr.currentSize = 0
	return nil
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.file.Close()
}

func (lr *LogRotator) ArchiveOldLogs(keepCount int) error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	pattern := lr.filePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= keepCount {
		return nil
	}

	oldestFiles := matches[:len(matches)-keepCount]
	for _, file := range oldestFiles {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
	}

	if err := rotator.ArchiveOldLogs(5); err != nil {
		fmt.Printf("Archive error: %v\n", err)
	}
}