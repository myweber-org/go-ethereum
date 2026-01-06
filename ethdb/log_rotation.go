
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
	maxFileSize = 1024 * 1024 * 10 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	baseName   string
	currentDay string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: filepath.Join(logDir, name),
	}
	if err := rl.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.rotateIfNeeded(); err != nil {
		return 0, err
	}

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}
	rl.size += int64(n)
	return n, nil
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	if rl.file == nil || rl.size >= maxFileSize || rl.currentDay != currentDate {
		if rl.file != nil {
			rl.file.Close()
		}

		rl.currentDay = currentDate
		filename := fmt.Sprintf("%s-%s.log", rl.baseName, currentDate)

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		info, err := file.Stat()
		if err != nil {
			return err
		}

		rl.file = file
		rl.size = info.Size()

		go rl.cleanupOldFiles()
	}
	return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
	files, err := filepath.Glob(rl.baseName + "-*.log")
	if err != nil {
		return
	}

	if len(files) > maxBackups {
		for i := 0; i < len(files)-maxBackups; i++ {
			os.Remove(files[i])
		}
	}
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d: Application is running smoothly", i)
		time.Sleep(100 * time.Millisecond)
	}
}