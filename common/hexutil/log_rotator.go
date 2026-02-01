
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

	rl := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	// Rotate existing backup files
	for i := rl.backupCount - 1; i > 0; i-- {
		oldName := fmt.Sprintf("%s.%d.gz", rl.filePath, i)
		newName := fmt.Sprintf("%s.%d.gz", rl.filePath, i+1)
		if _, err := os.Stat(oldName); err == nil {
			os.Rename(oldName, newName)
		}
	}

	// Compress current log file
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
	if err := os.Rename(rl.filePath, backupName); err != nil {
		return err
	}

	// Compress the backup file
	compressedName := fmt.Sprintf("%s.1.gz", rl.filePath)
	if err := compressFile(backupName, compressedName); err != nil {
		log.Printf("Failed to compress %s: %v", backupName, err)
		os.Remove(backupName)
	} else {
		os.Remove(backupName)
	}

	// Remove old backups beyond backupCount
	for i := rl.backupCount + 1; i <= rl.backupCount+10; i++ {
		oldFile := fmt.Sprintf("%s.%d.gz", rl.filePath, i)
		os.Remove(oldFile)
	}

	return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Simple compression simulation (in real implementation use compress/gzip)
	_, err = io.Copy(dstFile, srcFile)
	return err
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
	logger, err := NewRotatingLogger("./logs/app.log", 10, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Simulate log writing
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: %s\n",
			time.Now().Format(time.RFC3339),
			i,
			strings.Repeat("Test", i%10+1))
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}