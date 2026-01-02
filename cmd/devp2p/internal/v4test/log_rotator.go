package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize    = 10 * 1024 * 1024 // 10MB
	backupCount   = 5
	logFileName   = "app.log"
	checkInterval = 30 * time.Second
)

type LogRotator struct {
	currentSize int64
	file        *os.File
}

func NewLogRotator() (*LogRotator, error) {
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &LogRotator{
		currentSize: info.Size(),
		file:        file,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxLogSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.file.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if err := lr.file.Close(); err != nil {
		return err
	}

	// Rename current log file with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
	if err := os.Rename(logFileName, backupName); err != nil {
		return err
	}

	// Compress the backup
	if err := compressFile(backupName); err != nil {
		return err
	}

	// Clean old backups
	if err := cleanupOldBackups(); err != nil {
		return err
	}

	// Create new log file
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.file = file
	lr.currentSize = 0
	return nil
}

func compressFile(filename string) error {
	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}

	// Remove uncompressed file
	return os.Remove(filename)
}

func cleanupOldBackups() error {
	pattern := logFileName + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= backupCount {
		return nil
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(matches)-backupCount; i++ {
		if err := os.Remove(matches[i]); err != nil {
			return err
		}
	}

	return nil
}

func (lr *LogRotator) Close() error {
	return lr.file.Close()
}

func main() {
	rotator, err := NewLogRotator()
	if err != nil {
		log.Fatal(err)
	}
	defer rotator.Close()

	log.SetOutput(rotator)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Heartbeat at %s", time.Now().Format(time.RFC3339))
		}
	}
}