package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxFiles    = 5
	logFileName = "app.log"
)

type LogRotator struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewLogRotator(logDir string) (*LogRotator, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	basePath := filepath.Join(logDir, logFileName)
	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &LogRotator{
		currentFile: file,
		currentSize: info.Size(),
		basePath:    basePath,
	}, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
	if lr.currentSize+int64(len(p)) > maxFileSize {
		if err := lr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := lr.currentFile.Write(p)
	if err == nil {
		lr.currentSize += int64(n)
	}
	return n, err
}

func (lr *LogRotator) rotate() error {
	if err := lr.currentFile.Close(); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)
	if err := os.Rename(lr.basePath, rotatedPath); err != nil {
		return err
	}

	file, err := os.OpenFile(lr.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	lr.currentFile = file
	lr.currentSize = 0

	go lr.cleanupOldLogs()
	return nil
}

func (lr *LogRotator) cleanupOldLogs() {
	pattern := lr.basePath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) <= maxFiles {
		return
	}

	sort.Strings(matches)
	filesToRemove := matches[:len(matches)-maxFiles]

	for _, file := range filesToRemove {
		os.Remove(file)
	}
}

func (lr *LogRotator) Close() error {
	return lr.currentFile.Close()
}

func main() {
	rotator, err := NewLogRotator("./logs")
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	writer := io.MultiWriter(os.Stdout, rotator)
	for i := 0; i < 100; i++ {
		fmt.Fprintf(writer, "Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		time.Sleep(100 * time.Millisecond)
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

func (rl *RotatingLogger) rotate() error {
	if rl.currentFile != nil {
		rl.currentFile.Close()
		rl.compressCurrentFile()
	}

	rl.fileCount++
	if rl.fileCount > rl.maxFiles {
		rl.cleanupOldFiles()
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) openCurrentFile() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.log", rl.basePath, timestamp)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rl.currentFile = file
	rl.fileSize = info.Size()
	return nil
}

func (rl *RotatingLogger) compressCurrentFile() {
	if rl.currentFile == nil {
		return
	}

	oldPath := rl.currentFile.Name()
	newPath := oldPath + ".gz"

	oldFile, err := os.Open(oldPath)
	if err != nil {
		return
	}
	defer oldFile.Close()

	newFile, err := os.Create(newPath)
	if err != nil {
		return
	}
	defer newFile.Close()

	gzWriter := gzip.NewWriter(newFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, oldFile); err != nil {
		return
	}

	os.Remove(oldPath)
}

func (rl *RotatingLogger) cleanupOldFiles() {
	pattern := rl.basePath + "_*.log.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	if len(matches) > rl.maxFiles {
		filesToRemove := matches[:len(matches)-rl.maxFiles]
		for _, file := range filesToRemove {
			os.Remove(file)
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
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Rotator struct {
	mu          sync.Mutex
	file        *os.File
	currentSize int64
	maxSize     int64
	basePath    string
	rotateTime  time.Time
	interval    time.Duration
}

func NewRotator(basePath string, maxSize int64, interval time.Duration) (*Rotator, error) {
	if err := os.MkdirAll(filepath.Dir(basePath), 0755); err != nil {
		return nil, err
	}

	r := &Rotator{
		maxSize:  maxSize,
		basePath: basePath,
		interval: interval,
	}
	if err := r.openCurrent(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Rotator) openCurrent() error {
	now := time.Now()
	filename := fmt.Sprintf("%s.%s.log", r.basePath, now.Format("2006-01-02_150405"))
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	r.file = file
	r.currentSize = stat.Size()
	r.rotateTime = now
	return nil
}

func (r *Rotator) rotate() error {
	if r.file != nil {
		r.file.Close()
	}

	archiveName := fmt.Sprintf("%s.%s.log", r.basePath, time.Now().Format("2006-01-02_150405"))
	return os.Rename(r.basePath+".log", archiveName)
}

func (r *Rotator) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if r.currentSize+int64(len(p)) > r.maxSize || now.Sub(r.rotateTime) > r.interval {
		if err := r.rotate(); err != nil {
			return 0, err
		}
		if err := r.openCurrent(); err != nil {
			return 0, err
		}
	}

	n, err := r.file.Write(p)
	if err == nil {
		r.currentSize += int64(n)
	}
	return n, err
}

func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func main() {
	rotator, err := NewRotator("./logs/app", 10*1024*1024, 24*time.Hour)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create rotator: %v\n", err)
		os.Exit(1)
	}
	defer rotator.Close()

	for i := 0; i < 100; i++ {
		message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
		if _, err := rotator.Write([]byte(message)); err != nil {
			fmt.Fprintf(os.Stderr, "Write failed: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}