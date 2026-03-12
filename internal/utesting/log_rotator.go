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
	maxFileSize = 10 * 1024 * 1024
	maxBackups  = 5
	logDir      = "./logs"
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
	baseName    string
}

func NewRotatingLogger(name string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	rl := &RotatingLogger{
		baseName: name,
	}
	if err := rl.openCurrent(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > maxFileSize {
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
		timestamp := time.Now().Format("20060102_150405")
		oldPath := filepath.Join(logDir, rl.baseName+".log")
		newPath := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", rl.baseName, timestamp))
		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
		if err := rl.compressFile(newPath); err != nil {
			return err
		}
		rl.cleanupOld()
	}
	return rl.openCurrent()
}

func (rl *RotatingLogger) openCurrent() error {
	path := filepath.Join(logDir, rl.baseName+".log")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

	gz := gzip.NewWriter(destFile)
	defer gz.Close()

	if _, err := io.Copy(gz, srcFile); err != nil {
		return err
	}
	os.Remove(src)
	return nil
}

func (rl *RotatingLogger) cleanupOld() {
	pattern := filepath.Join(logDir, rl.baseName+"_*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}
	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, f := range toDelete {
			os.Remove(f)
		}
	}
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
	logger, err := NewRotatingLogger("app")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := logger.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	maxLogSize   = 1024 * 1024 // 1MB
	logFileName  = "app.log"
	archiveDir   = "archives"
)

func rotateLogIfNeeded() error {
	info, err := os.Stat(logFileName)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat log file: %w", err)
	}

	if info.Size() < maxLogSize {
		return nil
	}

	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	archiveName := filepath.Join(archiveDir, fmt.Sprintf("%s_%s", logFileName, timestamp))
	
	if err := os.Rename(logFileName, archiveName); err != nil {
		return fmt.Errorf("failed to archive log file: %w", err)
	}

	newFile, err := os.Create(logFileName)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}
	newFile.Close()

	fmt.Printf("Log rotated: %s -> %s\n", logFileName, archiveName)
	return nil
}

func writeLog(message string) error {
	if err := rotateLogIfNeeded(); err != nil {
		return err
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	logEntry := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), message)
	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}

	return nil
}

func main() {
	for i := 1; i <= 100; i++ {
		message := fmt.Sprintf("Log entry number %d", i)
		if err := writeLog(message); err != nil {
			fmt.Printf("Error writing log: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}