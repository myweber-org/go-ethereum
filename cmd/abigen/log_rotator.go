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
	maxFileSize    = 10 * 1024 * 1024
	maxBackupFiles = 5
	logDir         = "./logs"
)

type LogRotator struct {
	currentFile   *os.File
	currentSize   int64
	baseFilename  string
	fileExtension string
}

func NewLogRotator(filename string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filepath.Base(filename), ext)

	rotator := &LogRotator{
		baseFilename:  base,
		fileExtension: ext,
	}

	if err := rotator.openCurrentFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
	path := filepath.Join(logDir, lr.baseFilename+lr.fileExtension)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	backupName := fmt.Sprintf("%s_%s%s", lr.baseFilename, timestamp, lr.fileExtension)
	oldPath := filepath.Join(logDir, lr.baseFilename+lr.fileExtension)
	newPath := filepath.Join(logDir, backupName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := lr.openCurrentFile(); err != nil {
		return err
	}

	return lr.cleanupOldFiles()
}

func (lr *LogRotator) cleanupOldFiles() error {
	pattern := filepath.Join(logDir, lr.baseFilename+"_*"+lr.fileExtension)
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

func (lr *LogRotator) Close() error {
	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}