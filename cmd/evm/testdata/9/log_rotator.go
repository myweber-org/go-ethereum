
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024 // 10MB
	maxBackupCount = 5
	logExtension = ".log"
	compressExtension = ".gz"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewLogRotator(basePath string) (*LogRotator, error) {
	rotator := &LogRotator{
		basePath: strings.TrimSuffix(basePath, logExtension),
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if lr.currentFile != nil {
		lr.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := fmt.Sprintf("%s_%s%s", lr.basePath, timestamp, logExtension)
	if err := os.Rename(lr.currentFile.Name(), oldPath); err != nil {
		return err
	}

	if err := lr.compressFile(oldPath); err != nil {
		fmt.Printf("Warning: failed to compress %s: %v\n", oldPath, err)
	}

	if err := lr.cleanupOldFiles(); err != nil {
		fmt.Printf("Warning: failed to cleanup old files: %v\n", err)
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
	path := lr.basePath + logExtension
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = info.Size()
	return nil
}

func (lr *LogRotator) compressFile(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	destPath := path + compressExtension
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// Simple compression simulation (in real implementation use gzip)
	_, err = io.Copy(dest, src)
	if err != nil {
		os.Remove(destPath)
		return err
	}

	os.Remove(path)
	return nil
}

func (lr *LogRotator) cleanupOldFiles() error {
	pattern := lr.basePath + "_*" + logExtension + compressExtension
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackupCount {
		return nil
	}

	filesToRemove := matches[:len(matches)-maxBackupCount]
	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("application.log")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}