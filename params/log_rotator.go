package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	currentFile   *os.File
	currentSize   int64
	maxSize       int64
	logDirectory  string
	baseFilename  string
	retentionDays int
}

func NewRotatingLogger(dir, filename string, maxSize int64, retention int) (*RotatingLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	fullPath := filepath.Join(dir, filename)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &RotatingLogger{
		currentFile:   file,
		currentSize:   info.Size(),
		maxSize:       maxSize,
		logDirectory:  dir,
		baseFilename:  filename,
		retentionDays: retention,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, fmt.Errorf("rotation failed: %w", err)
		}
	}

	n, err := rl.currentFile.Write(p)
	if err != nil {
		return n, fmt.Errorf("write failed: %w", err)
	}

	rl.currentSize += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return fmt.Errorf("failed to close current file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(rl.logDirectory, rl.baseFilename)
	newPath := filepath.Join(rl.logDirectory, fmt.Sprintf("%s.%s.log", rl.baseFilename, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if err := rl.compressFile(newPath); err != nil {
		return fmt.Errorf("failed to compress old log: %w", err)
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	rl.currentFile = file
	rl.currentSize = 0

	go rl.cleanupOldFiles()

	return nil
}

func (rl *RotatingLogger) compressFile(source string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(source + ".gz")
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	if err := os.Remove(source); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
	if rl.retentionDays <= 0 {
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -rl.retentionDays)

	files, err := os.ReadDir(rl.logDirectory)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffTime) && (filepath.Ext(info.Name()) == ".gz" || filepath.Ext(info.Name()) == ".log") {
			oldPath := filepath.Join(rl.logDirectory, info.Name())
			os.Remove(oldPath)
		}
	}
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}