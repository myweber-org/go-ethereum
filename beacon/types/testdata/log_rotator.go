package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
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
	
	rotator := &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		rotationCount: 0,
	}
	
	if err := rotator.openFile(); err != nil {
		return nil, err
	}
	
	return rotator, nil
}

func (lr *LogRotator) openFile() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	
	if lr.file != nil {
		lr.file.Close()
	}
	
	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	
	lr.file = file
	lr.currentSize = info.Size()
	
	return nil
}

func (lr *LogRotator) rotate() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	
	if lr.file != nil {
		lr.file.Close()
	}
	
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s.%d", lr.filePath, timestamp, lr.rotationCount)
	
	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}
	
	lr.rotationCount++
	
	return lr.openFile()
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	
	if lr.currentSize+int64(len(p)) > lr.maxSize {
		lr.mu.Unlock()
		if err := lr.rotate(); err != nil {
			lr.mu.Lock()
			return 0, err
		}
		lr.mu.Lock()
	}
	
	n, err := lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	
	return n, err
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	
	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func (lr *LogRotator) CleanOldLogs(maxAgeDays int) error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	
	cutoffTime := time.Now().AddDate(0, 0, -maxAgeDays)
	
	files, err := filepath.Glob(lr.filePath + ".*")
	if err != nil {
		return err
	}
	
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		
		if info.ModTime().Before(cutoffTime) {
			os.Remove(file)
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
	
	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n", 
			time.Now().Format(time.RFC3339), i)
		_, err := rotator.Write([]byte(logEntry))
		if err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
		}
		
		time.Sleep(100 * time.Millisecond)
	}
	
	if err := rotator.CleanOldLogs(7); err != nil {
		fmt.Printf("Failed to clean old logs: %v\n", err)
	}
	
	fmt.Println("Log rotation test completed")
}