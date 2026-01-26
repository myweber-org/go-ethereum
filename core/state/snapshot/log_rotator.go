
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	writer     *bufio.Writer
	basePath   string
	currentNum int
	fileSize   int64
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
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
		rl.writer.Flush()
		rl.file.Close()
	}

	var err error
	rl.file, err = os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := rl.file.Stat()
	if err != nil {
		return err
	}

	rl.fileSize = info.Size()
	rl.writer = bufio.NewWriter(rl.file)

	rl.findCurrentRotationNumber()
	return nil
}

func (rl *RotatingLogger) findCurrentRotationNumber() {
	pattern := rl.basePath + ".*"
	matches, _ := filepath.Glob(pattern)
	maxNum := 0

	for _, match := range matches {
		ext := filepath.Ext(match)
		if num, err := strconv.Atoi(ext[1:]); err == nil && num > maxNum {
			maxNum = num
		}
	}

	rl.currentNum = maxNum
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.fileSize+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.writer.Write(p)
	if err != nil {
		return n, err
	}

	rl.fileSize += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	rl.writer.Flush()
	rl.file.Close()

	rl.currentNum++
	if rl.currentNum > maxBackups {
		rl.currentNum = 1
	}

	oldPath := rl.basePath
	newPath := fmt.Sprintf("%s.%d", rl.basePath, rl.currentNum)

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.openCurrentFile(); err != nil {
		return err
	}

	rl.cleanOldBackups()
	return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
	pattern := rl.basePath + ".*"
	matches, _ := filepath.Glob(pattern)

	backupFiles := make(map[int]string)
	for _, match := range matches {
		ext := filepath.Ext(match)
		if num, err := strconv.Atoi(ext[1:]); err == nil {
			backupFiles[num] = match
		}
	}

	for num, path := range backupFiles {
		if num > maxBackups {
			os.Remove(path)
		}
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.writer.Flush()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: This is a test log message.\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(message))
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}