package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotatingLogger struct {
	currentFile   *os.File
	maxSize       int64
	backupCount   int
	logDir        string
	baseName      string
	currentSize   int64
	compressOld   bool
}

func NewRotatingLogger(dir, name string, maxSizeMB int, backups int, compress bool) (*RotatingLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		maxSize:     maxSize,
		backupCount: backups,
		logDir:      dir,
		baseName:    name,
		compressOld: compress,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := filepath.Join(rl.logDir, rl.baseName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Failed to rotate log: %v", err)
		}
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(rl.logDir, rl.baseName)
	newPath := filepath.Join(rl.logDir, fmt.Sprintf("%s.%s", rl.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if rl.compressOld {
		go rl.compressFile(newPath)
	}

	if err := rl.openCurrentFile(); err != nil {
		return err
	}

	rl.cleanupOldFiles()
	return nil
}

func (rl *RotatingLogger) compressFile(path string) {
	// Compression implementation would go here
	// For simplicity, just rename with .gz extension
	compressedPath := path + ".gz"
	if err := os.Rename(path, compressedPath); err != nil {
		log.Printf("Failed to rename file for compression: %v", err)
	}
}

func (rl *RotatingLogger) cleanupOldFiles() {
	pattern := filepath.Join(rl.logDir, rl.baseName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= rl.backupCount {
		return
	}

	filesToDelete := matches[:len(matches)-rl.backupCount]
	for _, file := range filesToDelete {
		os.Remove(file)
		if rl.compressOld {
			os.Remove(file + ".gz")
		}
	}
}

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("./logs", "app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}
}