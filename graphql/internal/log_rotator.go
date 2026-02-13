
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type LogRotator struct {
	mu          sync.Mutex
	filePath    string
	maxSize     int64
	backupCount int
	currentFile *os.File
	currentSize int64
}

func NewLogRotator(filePath string, maxSize int64, backupCount int) (*LogRotator, error) {
	lr := &LogRotator{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := lr.openCurrentFile(); err != nil {
		return nil, err
	}

	return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	lr.currentFile = file
	lr.currentSize = stat.Size()
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	lr.currentSize += int64(n)
	return n, nil
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	for i := lr.backupCount - 1; i >= 0; i-- {
		oldPath := lr.getBackupPath(i)
		newPath := lr.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(lr.filePath, lr.getBackupPath(0)); err != nil {
		return err
	}

	return lr.openCurrentFile()
}

func (lr *LogRotator) getBackupPath(index int) string {
	if index == 0 {
		return lr.filePath + ".1"
	}
	return lr.filePath + "." + strconv.Itoa(index+1)
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentFile != nil {
		return lr.currentFile.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 1024*1024, 5)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: This is a sample log message.\n", i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed.")
}