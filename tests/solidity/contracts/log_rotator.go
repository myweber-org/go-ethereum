
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
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	mu        sync.Mutex
	file      *os.File
	size      int64
	baseName  string
	fileCount int
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, name)
	file, err := os.OpenFile(basePath+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	rl := &RotatingLogger{
		file:     file,
		size:     info.Size(),
		baseName: name,
	}

	rl.countExistingBackups()
	return rl, nil
}

func (rl *RotatingLogger) countExistingBackups() {
	pattern := filepath.Join(logDir, rl.baseName+".*.log.gz")
	matches, _ := filepath.Glob(pattern)
	rl.fileCount = len(matches)
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err := rl.file.Write(p)
	if err != nil {
		return n, err
	}

	rl.size += int64(n)
	if rl.size >= maxFileSize {
		if err := rl.rotate(); err != nil {
			log.Printf("Rotation failed: %v", err)
		}
	}
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	currentPath := filepath.Join(logDir, rl.baseName+".log")
	backupPath := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", rl.baseName, time.Now().Format("20060102150405")))

	if err := os.Rename(currentPath, backupPath); err != nil {
		return err
	}

	if err := rl.compressFile(backupPath); err != nil {
		log.Printf("Compression failed: %v", err)
	}

	file, err := os.OpenFile(currentPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.size = 0
	rl.fileCount++

	rl.cleanupOldBackups()
	return nil
}

func (rl *RotatingLogger) compressFile(src string) error {
	dest := src + ".gz"
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	gzWriter := gzip.NewWriter(destFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	os.Remove(src)
	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := filepath.Join(logDir, rl.baseName+".*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxBackups {
		return
	}

	oldest := matches[:len(matches)-maxBackups]
	for _, file := range oldest {
		os.Remove(file)
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("application")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	customLog := log.New(logger, "", log.LstdFlags)

	for i := 0; i < 1000; i++ {
		customLog.Printf("Log entry %d: Simulating application activity", i)
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation demonstration completed")
}