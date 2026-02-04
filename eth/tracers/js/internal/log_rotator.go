
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

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	// Find next available backup number
	rl.currentNum = 1
	for {
		backupPath := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.currentNum)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			break
		}
		rl.currentNum++
		if rl.currentNum > maxBackups {
			rl.removeOldestBackup()
			rl.currentNum = maxBackups
		}
	}

	// Compress current file
	source, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destPath := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.currentNum)
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, source); err != nil {
		return err
	}

	// Remove original and create new
	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) removeOldestBackup() {
	oldestNum := 1
	oldestTime := time.Now()

	for i := 1; i <= maxBackups; i++ {
		path := fmt.Sprintf("%s.%d.gz", rl.basePath, i)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().Before(oldestTime) {
			oldestTime = info.ModTime()
			oldestNum = i
		}
	}

	path := fmt.Sprintf("%s.%d.gz", rl.basePath, oldestNum)
	os.Remove(path)
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// Simulate log writing
	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: This is a test log message\n",
			time.Now().Format("2006-01-02 15:04:05"), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed. Check compressed backup files:")
	files, _ := filepath.Glob("app.log*")
	for _, f := range files {
		fmt.Println("  ", f)
	}
}