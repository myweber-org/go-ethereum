package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	file         *os.File
	basePath     string
	maxSize      int64
	maxFiles     int
	currentSize  int64
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	logger := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := logger.openOrCreate(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *RotatingLogger) openOrCreate() error {
	dir := filepath.Dir(l.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.OpenFile(l.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("stat log file: %w", err)
	}

	l.file = file
	l.currentSize = info.Size()
	return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentSize+int64(len(p)) > l.maxSize {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := l.file.Write(p)
	if err == nil {
		l.currentSize += int64(n)
	}
	return n, err
}

func (l *RotatingLogger) rotate() error {
	if err := l.file.Close(); err != nil {
		return fmt.Errorf("close current log: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", l.basePath, timestamp)

	if err := os.Rename(l.basePath, rotatedPath); err != nil {
		return fmt.Errorf("rename log file: %w", err)
	}

	if err := l.openOrCreate(); err != nil {
		return fmt.Errorf("reopen log file: %w", err)
	}

	l.cleanupOldFiles()
	return nil
}

func (l *RotatingLogger) cleanupOldFiles() {
	dir := filepath.Dir(l.basePath)
	baseName := filepath.Base(l.basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var rotatedFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, baseName+".") && entry.Type().IsRegular() {
			rotatedFiles = append(rotatedFiles, filepath.Join(dir, name))
		}
	}

	if len(rotatedFiles) <= l.maxFiles {
		return
	}

	sort.Strings(rotatedFiles)
	filesToRemove := rotatedFiles[:len(rotatedFiles)-l.maxFiles]

	for _, file := range filesToRemove {
		os.Remove(file)
	}
}

func (l *RotatingLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write failed: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}