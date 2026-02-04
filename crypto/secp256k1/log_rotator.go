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

type RotatingLogger struct {
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	backupCount  int
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024
	if maxSize <= 0 {
		return nil, fmt.Errorf("maxSize must be positive")
	}

	rl := &RotatingLogger{
		basePath:    basePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.basePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Failed to rotate log: %v", err)
		}
	}

	n, err = rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	for i := rl.backupCount - 1; i >= 0; i-- {
		oldPath := rl.getBackupPath(i)
		newPath := rl.getBackupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == rl.backupCount-1 {
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	if _, err := os.Stat(rl.basePath); err == nil {
		backupPath := rl.getBackupPath(0)
		os.Rename(rl.basePath, backupPath)
		rl.compressFile(backupPath)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) getBackupPath(index int) string {
	if index == 0 {
		return rl.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", rl.basePath, index)
}

func (rl *RotatingLogger) compressFile(path string) {
	go func() {
		src, err := os.Open(path)
		if err != nil {
			return
		}
		defer src.Close()

		dstPath := path + ".gz"
		dst, err := os.Create(dstPath)
		if err != nil {
			return
		}
		defer dst.Close()

		gzWriter := NewGzipWriter(dst)
		defer gzWriter.Close()

		io.Copy(gzWriter, src)
		os.Remove(path)
	}()
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: Application event occurred at %v", i, time.Now())
		time.Sleep(10 * time.Millisecond)
	}
}

type GzipWriter struct {
	io.WriteCloser
}

func NewGzipWriter(w io.Writer) *GzipWriter {
	return &GzipWriter{}
}

func (gw *GzipWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (gw *GzipWriter) Close() error {
	return nil
}