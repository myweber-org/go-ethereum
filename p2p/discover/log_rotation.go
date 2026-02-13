
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const maxLogSize = 1024 * 1024 // 1MB

type RotatingLogger struct {
	mu       sync.Mutex
	file     *os.File
	filePath string
	fileSize int64
	sequence int
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: basePath,
	}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.fileSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	if rl.fileSize < maxLogSize {
		return nil
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.sequence++
	backupPath := fmt.Sprintf("%s.%d", rl.filePath, rl.sequence)
	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	rl.file.Close()
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.file.Write(p)
	if err == nil {
		rl.fileSize += int64(n)
	}
	return n, err
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
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	for i := 0; i < 10000; i++ {
		message := fmt.Sprintf("Log entry %d: Application is running normally\n", i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
			break
		}
	}

	fmt.Println("Log rotation test completed")
}