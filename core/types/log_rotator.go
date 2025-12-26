package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
		fileIndex: 0,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}
	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	filename := fmt.Sprintf("%s.%d.log", rl.filePath, rl.fileIndex)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.currentSize >= rl.maxSize {
		rl.currentFile.Close()
		rl.fileIndex++
		return rl.openCurrentFile()
	}
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) WriteString(s string) (int, error) {
	return rl.Write([]byte(s))
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app", 1) // 1MB max size
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.WriteString(msg); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}