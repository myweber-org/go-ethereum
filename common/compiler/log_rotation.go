package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type RotatingLogger struct {
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
}

func NewRotatingLogger(path string, maxSize int64) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		filePath: path,
		maxSize:  maxSize,
	}

	if err := rl.openFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	info, err := os.Stat(rl.filePath)
	if err == nil {
		rl.currentSize = info.Size()
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	rl.file = file
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}

	rl.currentSize += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	if err := rl.openFile(); err != nil {
		return err
	}

	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) Close() error {
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}