package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const maxFileSize = 1024 * 1024 // 1 MB
const backupSuffix = ".bak"

type LogRotator struct {
	filePath string
	file     *os.File
	mu       sync.Mutex
}

func NewLogRotator(filePath string) (*LogRotator, error) {
	lr := &LogRotator{filePath: filePath}
	err := lr.openFile()
	return lr, err
}

func (lr *LogRotator) openFile() error {
	dir := filepath.Dir(lr.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	lr.file = file
	return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	info, err := lr.file.Stat()
	if err != nil {
		return 0, err
	}

	if info.Size()+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	return lr.file.Write(p)
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return err
	}

	backupPath := lr.filePath + backupSuffix
	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	return lr.openFile()
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator("logs/app.log")
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry number %d\n", i)
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}