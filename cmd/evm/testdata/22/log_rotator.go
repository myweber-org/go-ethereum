
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	backupCount = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentDay string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
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
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	today := time.Now().Format("2006-01-02")
	if rl.currentDay != today {
		if err := rl.rotateByDate(); err != nil {
			return err
		}
		rl.currentDay = today
	}

	if rl.size >= maxFileSize {
		if err := rl.rotateBySize(); err != nil {
			return err
		}
	}

	if rl.file == nil {
		if err := rl.openCurrentFile(); err != nil {
			return err
		}
	}
	return nil
}

func (rl *RotatingLogger) rotateByDate() error {
	if rl.file != nil {
		rl.file.Close()
		rl.file = nil
	}
	return rl.openCurrentFile()
}

func (rl *RotatingLogger) rotateBySize() error {
	if rl.file != nil {
		rl.file.Close()
		rl.file = nil
	}

	// Rename existing files
	for i := backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)
		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	// Move current to backup.0
	currentPath := rl.getCurrentPath()
	backupPath := rl.getBackupPath(0)
	if _, err := os.Stat(currentPath); err == nil {
		os.Rename(currentPath, backupPath)
	}

	// Compress old backup
	go rl.compressOldBackups()

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := rl.getCurrentPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

func (rl *RotatingLogger) getCurrentPath() string {
	return fmt.Sprintf("%s.%s.log", rl.basePath, time.Now().Format("2006-01-02"))
}

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return fmt.Sprintf("%s.%s.log.0", rl.basePath, time.Now().Format("2006-01-02"))
	}
	return fmt.Sprintf("%s.%s.log.%d.gz", rl.basePath, time.Now().Format("2006-01-02"), index)
}

func (rl *RotatingLogger) compressOldBackups() {
	for i := backupCount; i > 0; i-- {
		srcPath := rl.getBackupPath(i - 1)
		dstPath := rl.getBackupPath(i)
		if strings.HasSuffix(srcPath, ".gz") {
			continue
		}
		if _, err := os.Stat(srcPath); err == nil {
			if err := compressFile(srcPath, dstPath); err == nil {
				os.Remove(srcPath)
			}
		}
	}
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Simple compression simulation
	// In real implementation, use gzip.NewWriter
	_, err = io.Copy(dstFile, srcFile)
	return err
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
	logger, err := NewRotatingLogger("/var/log/myapp/application")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
	}
}