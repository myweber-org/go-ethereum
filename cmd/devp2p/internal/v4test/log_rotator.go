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
	logDir      = "./logs"
)

type RotatingFileHandler struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	mu          sync.Mutex
}

func NewRotatingFileHandler(filename string) (*RotatingFileHandler, error) {
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

	return &RotatingFileHandler{
		currentFile: file,
		currentSize: info.Size(),
		baseName:    filename,
	}, nil
}

func (h *RotatingFileHandler) Write(p []byte) (int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.currentSize+int64(len(p)) > maxFileSize {
		if err := h.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := h.currentFile.Write(p)
	if err == nil {
		h.currentSize += int64(n)
	}
	return n, err
}

func (h *RotatingFileHandler) rotate() error {
	if err := h.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(logDir, h.baseName)
	newPath := filepath.Join(logDir, fmt.Sprintf("%s.%s", h.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := h.compressFile(newPath); err != nil {
		return err
	}

	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	h.currentFile = file
	h.currentSize = 0
	h.cleanupOldFiles()
	return nil
}

func (h *RotatingFileHandler) compressFile(source string) error {
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

func (h *RotatingFileHandler) cleanupOldFiles() {
	pattern := filepath.Join(logDir, h.baseName+".*.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > backupCount {
		filesToRemove := matches[:len(matches)-backupCount]
		for _, file := range filesToRemove {
			os.Remove(file)
		}
	}
}

func (h *RotatingFileHandler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.currentFile.Close()
}

func main() {
	handler, err := NewRotatingFileHandler("application.log")
	if err != nil {
		panic(err)
	}
	defer handler.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		if _, err := handler.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
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
	filePath    string
	maxSize     int64
	backupCount int
	currentSize int64
}

func NewRotatingLogger(filePath string, maxSizeMB int, backupCount int) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rl := &RotatingLogger{
		filePath:    filePath,
		maxSize:     maxSize,
		backupCount: backupCount,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	dir := filepath.Dir(rl.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat log file: %w", err)
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
			return 0, err
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

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)

	if err := os.Rename(rl.filePath, backupPath); err != nil {
		return fmt.Errorf("failed to rename log file: %w", err)
	}

	if err := rl.compressFile(backupPath); err != nil {
		return fmt.Errorf("failed to compress log file: %w", err)
	}

	if err := rl.cleanupOldBackups(); err != nil {
		return fmt.Errorf("failed to cleanup old backups: %w", err)
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(sourcePath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	compressedPath := sourcePath + ".gz"
	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		return err
	}
	defer compressedFile.Close()

	gzWriter := gzip.NewWriter(compressedFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, sourceFile); err != nil {
		return err
	}

	os.Remove(sourcePath)
	return nil
}

func (rl *RotatingLogger) cleanupOldBackups() error {
	pattern := rl.filePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) <= rl.backupCount {
		return nil
	}

	backups := make([]string, len(matches))
	copy(backups, matches)

	for i := 0; i < len(backups)-rl.backupCount; i++ {
		if err := os.Remove(backups[i]); err != nil {
			return err
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