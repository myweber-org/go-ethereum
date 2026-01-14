package main

import (
	"compress/gzip"
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
	backupCount int
	currentSize int64
}

func NewRotatingLogger(filePath string, maxSizeMB int, backupCount int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	logger := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(l.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	l.currentFile = file
	l.currentSize = stat.Size()
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

	n, err := l.currentFile.Write(p)
	if err == nil {
		l.currentSize += int64(n)
	}
	return n, err
}

func (l *RotatingLogger) rotate() error {
	if l.currentFile != nil {
		l.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", l.filePath, timestamp)

	if err := os.Rename(l.filePath, backupPath); err != nil {
		return err
	}

	if err := l.compressFile(backupPath); err != nil {
		return err
	}

	if err := l.cleanupOldBackups(); err != nil {
		return err
	}

	return l.openCurrentFile()
}

func (l *RotatingLogger) compressFile(sourcePath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	compressedPath := sourcePath + ".gz"
	destFile, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, sourceFile); err != nil {
		return err
	}

	os.Remove(sourcePath)
	return nil
}

func (l *RotatingLogger) cleanupOldBackups() error {
	pattern := l.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= l.backupCount {
		return nil
	}

	backupsToDelete := matches[:len(matches)-l.backupCount]
	for _, backup := range backupsToDelete {
		os.Remove(backup)
	}

	return nil
}

func (l *RotatingLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentFile != nil {
		return l.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10, 5)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}