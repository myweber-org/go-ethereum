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
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	baseName  string
	fileIndex int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
	}

	path := filepath.Join(logDir, rl.baseName+".log")
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
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file == nil {
		return fmt.Errorf("no open file")
	}

	rl.file.Close()
	oldPath := filepath.Join(logDir, rl.baseName+".log")
	newPath := filepath.Join(logDir, fmt.Sprintf("%s.%d.log", rl.baseName, rl.fileIndex))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	rl.fileIndex++
	if rl.fileIndex > maxBackups {
		rl.fileIndex = 0
	}

	if err := rl.cleanupOld(); err != nil {
		log.Printf("cleanup error: %v", err)
	}

	return rl.openCurrent()
}

func (rl *RotatingLogger) cleanupOld() error {
	pattern := filepath.Join(logDir, rl.baseName+".*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		for i := 0; i < len(matches)-maxBackups; i++ {
			if err := os.Remove(matches[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file == nil {
		return 0, fmt.Errorf("logger closed")
	}

	n, err = rl.file.Write(p)
	if err != nil {
		return n, err
	}
	rl.size += int64(n)

	if rl.size >= maxFileSize {
		go func() {
			if err := rl.rotate(); err != nil {
				log.Printf("rotate failed: %v", err)
			}
		}()
	}
	return n, nil
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

	multiWriter := io.MultiWriter(os.Stdout, logger)
	log.SetOutput(multiWriter)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(10 * time.Millisecond)
	}
}