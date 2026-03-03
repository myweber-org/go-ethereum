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
	fileCount   int
	maxFiles    int
	currentSize int64
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
	rl.currentSize = info.Size()
	return nil
}

func (rl *RotatingLogger) rotate() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentFile == nil {
		return fmt.Errorf("no current file")
	}

	rl.currentFile.Close()
	timestamp := time.Now().Format("20060102_150405")
	archivePath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

	if err := rl.compressFile(rl.basePath, archivePath); err != nil {
		return err
	}

	if err := os.Remove(rl.basePath); err != nil {
		return err
	}

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		rl.cleanOldFiles()
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	_, err = io.Copy(gz, source)
	return err
}

func (rl *RotatingLogger) cleanOldFiles() {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxFiles {
		filesToDelete := matches[:len(matches)-rl.maxFiles]
		for _, f := range filesToDelete {
			os.Remove(f)
		}
	}
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		rl.mu.Unlock()
		if err := rl.rotate(); err != nil {
			return 0, err
		}
		rl.mu.Lock()
	}

	n, err := rl.currentFile.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
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
}package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

func rotateLog(logFilePath string) error {
    // Check if the log file exists
    if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
        return fmt.Errorf("log file does not exist: %s", logFilePath)
    }

    // Generate timestamp suffix
    timestamp := time.Now().Format("20060102_150405")
    ext := filepath.Ext(logFilePath)
    baseName := logFilePath[:len(logFilePath)-len(ext)]
    rotatedFilePath := fmt.Sprintf("%s_%s%s", baseName, timestamp, ext)

    // Rename the current log file
    err := os.Rename(logFilePath, rotatedFilePath)
    if err != nil {
        return fmt.Errorf("failed to rename log file: %v", err)
    }

    // Create a new empty log file
    newFile, err := os.Create(logFilePath)
    if err != nil {
        // Attempt to rollback the rename if creating new file fails
        os.Rename(rotatedFilePath, logFilePath)
        return fmt.Errorf("failed to create new log file: %v", err)
    }
    newFile.Close()

    fmt.Printf("Log rotated successfully. Old file: %s\n", rotatedFilePath)
    return nil
}

func main() {
    // Example usage: rotate a log file named "application.log"
    err := rotateLog("application.log")
    if err != nil {
        fmt.Printf("Error rotating log: %v\n", err)
        os.Exit(1)
    }
}
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
	maxFileSize    = 10 * 1024 * 1024 // 10MB
	backupCount    = 5
	checkInterval  = 30 * time.Second
	compressBackup = true
)

type RotatingLogger struct {
	mu         sync.Mutex
	file       *os.File
	filePath   string
	currentSize int64
	baseName   string
	dir        string
}

func NewRotatingLogger(filePath string) (*RotatingLogger, error) {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	rl := &RotatingLogger{
		file:       file,
		filePath:   filePath,
		currentSize: info.Size(),
		baseName:   base,
		dir:        dir,
	}

	go rl.monitor()
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	n, err = rl.file.Write(p)
	if err != nil {
		return n, err
	}
	rl.currentSize += int64(n)

	if rl.currentSize >= maxFileSize {
		if err := rl.rotate(); err != nil {
			log.Printf("rotation failed: %v", err)
		}
	}
	return n, nil
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return fmt.Errorf("failed to close current file: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(rl.dir, fmt.Sprintf("%s.%s", rl.baseName, timestamp))

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename current file: %w", err)
	}

	if compressBackup {
		go compressFile(backupPath)
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	rl.file = file
	rl.currentSize = 0
	rl.cleanupOldBackups()
	return nil
}

func compressFile(path string) {
	// Compression implementation would go here
	// For now, just log the compression attempt
	log.Printf("Compressing backup file: %s", path)
	// In real implementation, use compress/gzip or similar
}

func (rl *RotatingLogger) cleanupOldBackups() {
	pattern := filepath.Join(rl.dir, rl.baseName+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("failed to list backup files: %v", err)
		return
	}

	if len(matches) <= backupCount {
		return
	}

	// Sort by modification time (oldest first)
	// For simplicity, we'll just remove the first N-oldest
	// In production, implement proper time-based sorting
	for i := 0; i < len(matches)-backupCount; i++ {
		if err := os.Remove(matches[i]); err != nil {
			log.Printf("failed to remove old backup %s: %v", matches[i], err)
		}
	}
}

func (rl *RotatingLogger) monitor() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		if rl.currentSize >= maxFileSize {
			if err := rl.rotate(); err != nil {
				log.Printf("scheduled rotation failed: %v", err)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func main() {
	logger, err := NewRotatingLogger("./logs/application.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	// Example usage
	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry number %d\n", 
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			log.Printf("write error: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}