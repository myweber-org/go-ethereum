
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogRotator struct {
	filePath    string
	maxSize     int64
	currentSize int64
	file        *os.File
	mu          sync.Mutex
}

func NewLogRotator(filePath string, maxSize int64) (*LogRotator, error) {
	rotator := &LogRotator{
		filePath: filePath,
		maxSize:  maxSize,
	}

	if err := rotator.openFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (lr *LogRotator) openFile() error {
	info, err := os.Stat(lr.filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if info != nil {
		lr.currentSize = info.Size()
	}

	file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.file = file
	return nil
}

func (lr *LogRotator) rotate() error {
	if lr.file != nil {
		lr.file.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

	if err := os.Rename(lr.filePath, backupPath); err != nil {
		return err
	}

	lr.currentSize = 0
	return lr.openFile()
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.currentSize+int64(len(p)) > lr.maxSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}

	return n, err
}

func (lr *LogRotator) Close() error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if lr.file != nil {
		return lr.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewLogRotator("app.log", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create log rotator: %v\n", err)
		return
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		rotator.Write([]byte(message))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}package main

import (
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
	currentSize int64
	rotationNum int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	dir := filepath.Dir(basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	logger := &RotatingLogger{
		filePath: basePath,
		maxSize:  maxSize,
	}

	if err := logger.openCurrentFile(); err != nil {
		return nil, err
	}

	return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.currentFile.Close()

	backupPath := fmt.Sprintf("%s.%d.%s", rl.filePath, rl.rotationNum, time.Now().Format("20060102_150405"))
	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return err
	}

	rl.rotationNum++
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

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.currentFile.Close()
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}