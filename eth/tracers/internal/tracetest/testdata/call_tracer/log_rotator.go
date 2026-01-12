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
}

func NewRotatingLogger(filePath string, maxSize int64, backupCount int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	rl.currentFile = file
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	info, err := rl.currentFile.Stat()
	if err != nil {
		return 0, fmt.Errorf("stat file: %w", err)
	}

	if info.Size()+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, fmt.Errorf("rotate logs: %w", err)
		}
	}

	return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return fmt.Errorf("close current file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return fmt.Errorf("rename file: %w", err)
	}

	if err := rl.compressBackup(backupPath); err != nil {
		return fmt.Errorf("compress backup: %w", err)
	}

	if err := rl.cleanOldBackups(); err != nil {
		return fmt.Errorf("clean old backups: %w", err)
	}

	return rl.openFile()
}

func (rl *RotatingLogger) compressBackup(srcPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	dstPath := srcPath + ".gz"
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return fmt.Errorf("compress data: %w", err)
	}

	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("remove original: %w", err)
	}

	return nil
}

func (rl *RotatingLogger) cleanOldBackups() error {
	pattern := rl.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob pattern: %w", err)
	}

	if len(matches) <= rl.backupCount {
		return nil
	}

	toDelete := matches[:len(matches)-rl.backupCount]
	for _, path := range toDelete {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("remove old backup %s: %w", path, err)
		}
	}

	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}