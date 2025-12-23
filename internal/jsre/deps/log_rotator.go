
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	rl.size += int64(n)
	return n, err
}

func (rl *RotatingLogger) rotateIfNeeded() error {
	now := time.Now()
	dateStr := now.Format("2006-01-02")

	if rl.file == nil || rl.currentDay != dateStr || rl.size >= maxFileSize {
		return rl.performRotation(dateStr)
	}
	return nil
}

func (rl *RotatingLogger) performRotation(dateStr string) error {
	if rl.file != nil {
		rl.file.Close()
		if err := rl.compressOldFile(); err != nil {
			return err
		}
	}

	newPath := fmt.Sprintf("%s.%s.log", rl.basePath, dateStr)
	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
	rl.currentDay = dateStr
	rl.cleanupOldBackups()
	return nil
}

func (rl *RotatingLogger) compressOldFile() error {
	if rl.currentDay == "" {
		return nil
	}

	oldPath := fmt.Sprintf("%s.%s.log", rl.basePath, rl.currentDay)
	compressedPath := oldPath + ".gz"

	src, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err = io.Copy(gz, src); err != nil {
		return err
	}

	return os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := rl.basePath + ".*.log.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > backupCount {
		sortByModTime(matches)
		for i := 0; i < len(matches)-backupCount; i++ {
			os.Remove(matches[i])
		}
	}
}

func sortByModTime(files []string) {
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			infoI, _ := os.Stat(files[i])
			infoJ, _ := os.Stat(files[j])
			if infoI.ModTime().After(infoJ.ModTime()) {
				files[i], files[j] = files[j], files[i]
			}
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
	logger, err := NewRotatingLogger("/var/log/myapp")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		logger.Write([]byte(fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))))
		time.Sleep(100 * time.Millisecond)
	}
}