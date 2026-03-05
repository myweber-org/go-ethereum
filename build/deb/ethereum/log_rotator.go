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
	mu           sync.Mutex
	currentFile  *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	rotationCount int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: basePath,
		maxSize:  int64(maxSizeMB) * 1024 * 1024,
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
	
	file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}
	
	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
	}
	
	rl.rotationCount++
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
	
	gzWriter := gzip.NewWriter(dest)
	defer gzWriter.Close()
	
	if _, err := io.Copy(gzWriter, source); err != nil {
		return err
	}
	
	if err := os.Remove(rl.basePath); err != nil {
		return err
	}
	
	return rl.openCurrentFile()
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
	logger, err := NewRotatingLogger("app.log", 10)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()
	
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	
	fmt.Println("Log rotation test completed")
}package main

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
	maxBackups  = 5
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	size       int64
	basePath   string
	currentNum int
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{
		basePath: path,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.file != nil {
		rl.file.Close()
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

	rl.file = file
	rl.size = info.Size()
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.size+int64(len(p)) > maxFileSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := rl.file.Write(p)
	if err == nil {
		rl.size += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	// Find next available rotation number
	for i := 1; i <= maxBackups; i++ {
		compressedPath := fmt.Sprintf("%s.%d.gz", rl.basePath, i)
		if _, err := os.Stat(compressedPath); os.IsNotExist(err) {
			rl.currentNum = i
			break
		}
	}

	// Compress current file
	sourcePath := rl.basePath
	destPath := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.currentNum)

	if err := compressFile(sourcePath, destPath); err != nil {
		return err
	}

	// Remove original file after compression
	if err := os.Remove(sourcePath); err != nil {
		return err
	}

	// Clean up old backups if exceeding maxBackups
	rl.cleanOldBackups()

	// Open new log file
	return rl.openCurrent()
}

func compressFile(source, dest string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()

	_, err = io.Copy(gz, in)
	return err
}

func (rl *RotatingLogger) cleanOldBackups() {
	for i := rl.currentNum + 1; i <= maxBackups; i++ {
		path := fmt.Sprintf("%s.%d.gz", rl.basePath, i)
		os.Remove(path)
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
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}