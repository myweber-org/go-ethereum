package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize  = 10 * 1024 * 1024
	maxBackups   = 5
	logDirectory = "./logs"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	basePath  string
	sequence  int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDirectory, 0755); err != nil {
		return nil, err
	}
	basePath := filepath.Join(logDirectory, baseName)
	rl := &RotatingLogger{basePath: basePath}
	if err := rl.openFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openFile() error {
	path := rl.basePath
	if rl.sequence > 0 {
		path = fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.file = file
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.file.Close()
	rl.sequence++
	if rl.sequence > maxBackups {
		rl.sequence = maxBackups
		oldest := fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)
		if err := os.Remove(oldest); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	for i := rl.sequence; i > 1; i-- {
		oldName := fmt.Sprintf("%s.%d", rl.basePath, i-1)
		newName := fmt.Sprintf("%s.%d", rl.basePath, i)
		if err := os.Rename(oldName, newName); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.Rename(rl.basePath, rl.basePath+".1"); err != nil {
		return err
	}
	rl.sequence = 0
	return rl.openFile()
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = rl.file.Write(p)
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logger))
	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}
}