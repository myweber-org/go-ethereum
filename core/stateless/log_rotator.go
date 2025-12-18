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

type RotatingLogger struct {
	mu          sync.Mutex
	currentFile *os.File
	basePath    string
	maxSize     int64
	fileSize    int64
	fileCount   int
	maxFiles    int
}

func NewRotatingLogger(basePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	f, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	rl.currentFile = f
	rl.fileSize = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.fileSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.fileSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

	if err := os.Rename(rl.basePath, rotatedPath); err != nil {
		return err
	}

	if err := rl.compressFile(rotatedPath); err != nil {
		return err
	}

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		if err := rl.cleanOldFiles(); err != nil {
			return err
		}
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(srcPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstPath := srcPath + ".gz"
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	if err := os.Remove(srcPath); err != nil {
		return err
	}

	return nil
}

func (rl *RotatingLogger) cleanOldFiles() error {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > rl.maxFiles {
		filesToDelete := matches[:len(matches)-rl.maxFiles]
		for _, file := range filesToDelete {
			if err := os.Remove(file); err != nil {
				return err
			}
		}
	}

	return nil
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
	logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}