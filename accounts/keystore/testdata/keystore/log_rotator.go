
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

func (rl *RotatingLogger) rotate() error {
	rl.currentFile.Close()

	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

	source, err := os.Open(rl.basePath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	if _, err := io.Copy(gz, source); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		rl.cleanupOldFiles()
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) cleanupOldFiles() {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxFiles {
		filesToDelete := len(matches) - rl.maxFiles
		for i := 0; i < filesToDelete; i++ {
			os.Remove(matches[i])
		}
	}
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

	for i := 0; i < 1000; i++ {
		logger.Write([]byte(fmt.Sprintf("Log entry %d: %s\n", i, time.Now().String())))
		time.Sleep(10 * time.Millisecond)
	}
}