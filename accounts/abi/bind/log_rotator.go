
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxFileSize   = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logDir        = "./logs"
	logPrefix     = "app"
	logExtension  = ".log"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewLogRotator() (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, logPrefix)
	rotator := &LogRotator{basePath: basePath}

	if err := rotator.openCurrentLog(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentLog() error {
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
	newPath := fmt.Sprintf("%s_%s%s", lr.basePath, timestamp, logExtension)
	oldPath := lr.basePath + logExtension

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := lr.openCurrentLog(); err != nil {
		return err
	}

	return lr.cleanupOldLogs()
}

func (lr *LogRotator) cleanupOldLogs() error {
	pattern := lr.basePath + "_*" + logExtension
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= maxBackupFiles {
		return nil
	}

	sort.Strings(matches)
	filesToRemove := matches[:len(matches)-maxBackupFiles]

	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) parseTimestamp(filename string) (time.Time, error) {
	base := strings.TrimSuffix(filepath.Base(filename), logExtension)
	parts := strings.Split(base, "_")
	if len(parts) < 3 {
		return time.Time{}, fmt.Errorf("invalid filename format")
	}

	timestampStr := parts[len(parts)-2] + "_" + parts[len(parts)-1]
	return time.Parse("20060102_150405", timestampStr)
}

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator()
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}

		if i%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("Log rotation test completed")
}