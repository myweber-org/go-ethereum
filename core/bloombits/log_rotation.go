package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RotatingLogger struct {
	mu         sync.Mutex
	basePath   string
	maxSize    int64
	maxBackups int
	current    *os.File
	size       int64
}

func NewRotatingLogger(basePath string, maxSize int64, maxBackups int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath:   basePath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
	}

	if err := rl.openOrCreate(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openOrCreate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.current != nil {
		rl.current.Close()
	}

	f, err := os.OpenFile(rl.basePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	rl.current = f
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) >= rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.current.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.current.Close(); err != nil {
		return err
	}

	for i := rl.maxBackups - 1; i >= 0; i-- {
		oldPath := rl.backupPath(i)
		newPath := rl.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(rl.basePath, rl.backupPath(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	return rl.openOrCreate()
}

func (rl *RotatingLogger) backupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".0"
	}
	return rl.basePath + "." + strconv.Itoa(index)
}

func (rl *RotatingLogger) Cleanup() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for i := rl.maxBackups; ; i++ {
		path := rl.backupPath(i)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			break
		}
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.current != nil {
		return rl.current.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 1024*10, 3)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	if err := logger.Cleanup(); err != nil {
		fmt.Printf("Cleanup error: %v\n", err)
	}

	fmt.Println("Log rotation test completed")
}