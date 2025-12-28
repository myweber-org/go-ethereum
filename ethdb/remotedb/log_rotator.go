
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	filePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	rl := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile == nil {
		return fmt.Errorf("no current log file")
	}

	rl.currentFile.Close()
	rl.rotationCount++

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%d.%s.gz", rl.filePath, rl.rotationCount, timestamp)

	if err := rl.compressFile(rl.filePath, archivePath); err != nil {
		return fmt.Errorf("compression failed: %v", err)
	}

	if err := os.Truncate(rl.filePath, 0); err != nil {
		return fmt.Errorf("truncate failed: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(source, target string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile == nil {
		return 0, fmt.Errorf("log file not open")
	}

	n, err := rl.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentSize += int64(n)
	if rl.currentSize >= rl.maxSize {
		go func() {
			if err := rl.rotate(); err != nil {
				log.Printf("Log rotation failed: %v", err)
			}
		}()
	}

	return n, nil
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
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal(err)
	}

	logPath := filepath.Join(logDir, "application.log")
	rotator, err := NewRotatingLogger(logPath, 10)
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(rotator)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation demonstration completed")
}