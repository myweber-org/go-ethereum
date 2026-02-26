package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Rotator struct {
	filePath    string
	maxSize     int64
	maxFiles    int
	currentSize int64
	file        *os.File
}

func NewRotator(filePath string, maxSizeMB int, maxFiles int) (*Rotator, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	rotator := &Rotator{
		filePath: filePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := rotator.openFile(); err != nil {
		return nil, err
	}

	return rotator, nil
}

func (r *Rotator) openFile() error {
	info, err := os.Stat(r.filePath)
	if err == nil {
		r.currentSize = info.Size()
	} else if os.IsNotExist(err) {
		r.currentSize = 0
	} else {
		return err
	}

	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	r.file = file
	return nil
}

func (r *Rotator) Write(p []byte) (int, error) {
	if r.currentSize+int64(len(p)) > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := r.file.Write(p)
	if err != nil {
		return n, err
	}

	r.currentSize += int64(n)
	return n, nil
}

func (r *Rotator) rotate() error {
	if err := r.file.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedFile := fmt.Sprintf("%s.%s", r.filePath, timestamp)

	if err := os.Rename(r.filePath, rotatedFile); err != nil {
		return err
	}

	if err := r.cleanupOldFiles(); err != nil {
		fmt.Printf("Warning: cleanup failed: %v\n", err)
	}

	r.currentSize = 0
	return r.openFile()
}

func (r *Rotator) cleanupOldFiles() error {
	dir := filepath.Dir(r.filePath)
	base := filepath.Base(r.filePath)

	files, err := filepath.Glob(filepath.Join(dir, base+".*"))
	if err != nil {
		return err
	}

	if len(files) <= r.maxFiles {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return extractTimestamp(files[i]) > extractTimestamp(files[j])
	})

	for i := r.maxFiles; i < len(files); i++ {
		if err := os.Remove(files[i]); err != nil {
			return err
		}
	}

	return nil
}

func extractTimestamp(filename string) time.Time {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return time.Time{}
	}

	timestamp := parts[len(parts)-1]
	t, err := time.Parse("20060102_150405", timestamp)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (r *Rotator) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewRotator("app.log", 10, 5)
	if err != nil {
		fmt.Printf("Failed to create rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 1000; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here.\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := rotator.Write([]byte(logEntry)); err != nil {
			fmt.Printf("Write failed: %v\n", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
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

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	baseName    string
	mu          sync.Mutex
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	rl := &RotatingLogger{
		baseName: baseName,
	}

	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}

	return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	path := filepath.Join(logDir, rl.baseName+".log")
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
	}

	timestamp := time.Now().Format("20060102_150405")
	oldPath := filepath.Join(logDir, rl.baseName+".log")
	newPath := filepath.Join(logDir, fmt.Sprintf("%s.%s.log", rl.baseName, timestamp))

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	if err := rl.compressFile(newPath); err != nil {
		fmt.Printf("Failed to compress %s: %v\n", newPath, err)
	}

	if err := rl.cleanupOldFiles(); err != nil {
		fmt.Printf("Failed to cleanup old files: %v\n", err)
	}

	return rl.openCurrentFile()
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

func (rl *RotatingLogger) cleanupOldFiles() error {
	pattern := filepath.Join(logDir, rl.baseName+".*.log.gz")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > backupCount {
		filesToRemove := matches[:len(matches)-backupCount]
		for _, file := range filesToRemove {
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
	logger, err := NewRotatingLogger("application")
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		message := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := logger.Write([]byte(message)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}