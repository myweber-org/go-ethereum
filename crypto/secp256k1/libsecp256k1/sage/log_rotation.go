package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logFileName    = "app.log"
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	rl.currentSize += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(rl.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))
	if err := os.Rename(filepath.Join(rl.basePath, logFileName), backupPath); err != nil {
		return err
	}

	if err := rl.cleanOldBackups(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
	fullPath := filepath.Join(rl.basePath, logFileName)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) cleanOldBackups() error {
	pattern := filepath.Join(rl.basePath, logFileName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackupFiles {
		return nil
	}

	oldestFirst := matches[:len(matches)-maxBackupFiles]
	for _, file := range oldestFirst {
		if err := os.Remove(file); err != nil {
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
	logger, err := NewRotatingLogger(".")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	writer := io.MultiWriter(os.Stdout, logger)
	for i := 0; i < 100; i++ {
		fmt.Fprintf(writer, "Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}