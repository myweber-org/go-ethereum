
package main

import (
	"compress/gzip"
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
	backupCount = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	filename   string
	current    *os.File
	size       int64
	mu         sync.Mutex
	rotateChan chan struct{}
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(logDir, filename)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	rl := &RotatingLogger{
		filename:   fullPath,
		current:    file,
		size:       info.Size(),
		rotateChan: make(chan struct{}, 1),
	}

	go rl.monitorRotation()
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err = rl.current.Write(p)
	if err != nil {
		return n, err
	}

	rl.size += int64(n)
	if rl.size >= maxFileSize {
		select {
		case rl.rotateChan <- struct{}{}:
		default:
		}
	}
	return n, nil
}

func (rl *RotatingLogger) monitorRotation() {
	for range rl.rotateChan {
		if err := rl.performRotation(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}
}

func (rl *RotatingLogger) performRotation() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if err := rl.current.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s.%s", rl.filename, timestamp)
	if err := os.Rename(rl.filename, backupName); err != nil {
		return err
	}

	if err := rl.compressFile(backupName); err != nil {
		log.Printf("Compression failed for %s: %v", backupName, err)
	}

	file, err := os.OpenFile(rl.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.current = file
	rl.size = 0
	rl.cleanupOldBackups()
	return nil
}

func (rl *RotatingLogger) compressFile(source string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(source + ".gz")
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	if err := os.Remove(source); err != nil {
		return err
	}
	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := rl.filename + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= backupCount {
		return
	}

	for i := 0; i < len(matches)-backupCount; i++ {
		os.Remove(matches[i])
	}
}

func (rl *RotatingLogger) Close() error {
	close(rl.rotateChan)
	return rl.current.Close()
}

func main() {
	logger, err := NewRotatingLogger("application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: This is a sample log message for testing rotation", i)
		time.Sleep(10 * time.Millisecond)
	}
}