
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
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	maxBackupFiles = 5
	logFileName    = "app.log"
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	currentPos int64
	basePath   string
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
	rl := &RotatingLogger{basePath: basePath}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	fullPath := filepath.Join(rl.basePath, logFileName)
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	rl.file = file
	rl.currentPos = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentPos+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}
	n, err = rl.file.Write(p)
	if err == nil {
		rl.currentPos += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	oldPath := filepath.Join(rl.basePath, logFileName)
	for i := maxBackupFiles - 1; i >= 0; i-- {
		var source string
		if i == 0 {
			source = oldPath
		} else {
			source = filepath.Join(rl.basePath, fmt.Sprintf("%s.%d", logFileName, i))
		}
		dest := filepath.Join(rl.basePath, fmt.Sprintf("%s.%d", logFileName, i+1))

		if _, err := os.Stat(source); err == nil {
			if err := os.Rename(source, dest); err != nil {
				return err
			}
		}
	}

	if err := rl.cleanupOldFiles(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	for i := maxBackupFiles + 1; i < 20; i++ {
		path := filepath.Join(rl.basePath, fmt.Sprintf("%s.%d", logFileName, i))
		if _, err := os.Stat(path); err == nil {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal(err)
	}

	rotator, err := NewRotatingLogger(logDir)
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, rotator))

	for i := 0; i < 100; i++ {
		log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}