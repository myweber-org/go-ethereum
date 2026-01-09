package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu          sync.Mutex
	basePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
	sequence    int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	logger := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
		sequence: 0,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}
	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := rl.basePath
	if rl.sequence > 0 {
		filename = fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.file = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.sequence++
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
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
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Some sample log data here.\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed.")
}