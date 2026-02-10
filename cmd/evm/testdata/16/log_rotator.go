package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize   = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logFileName   = "app.log"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{basePath: path}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	fullPath := filepath.Join(rl.basePath, logFileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	if rl.currentSize+int64(len(p)) > maxFileSize {
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
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(rl.basePath, logFileName)
	newPath := filepath.Join(rl.basePath, fmt.Sprintf("app_%s.log", timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.cleanupOldFiles(); err != nil {
		log.Printf("Failed to cleanup old log files: %v", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	pattern := filepath.Join(rl.basePath, "app_*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackupFiles {
		return nil
	}

	for i := 0; i < len(matches)-maxBackupFiles; i++ {
		if err := os.Remove(matches[i]); err != nil {
			return err
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
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

	rotator, err := NewRotatingLogger(logDir)
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, rotator))

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry number %d: Application is running normally", i)
		time.Sleep(10 * time.Millisecond)
	}
}