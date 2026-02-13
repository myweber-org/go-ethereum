
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

type RotatingLogger struct {
	mu           sync.Mutex
	file         *os.File
	basePath     string
	maxSize      int64
	currentSize  int64
	backupCount  int
	compressOld  bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, backupCount int, compressOld bool) (*RotatingLogger, error) {
	maxSize := int64(maxSizeMB) * 1024 * 1024

	file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &RotatingLogger{
		file:        file,
		basePath:    basePath,
		maxSize:     maxSize,
		currentSize: info.Size(),
		backupCount: backupCount,
		compressOld: compressOld,
	}, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.currentSize+int64(len(p)) > rl.maxSize {
		if err := rl.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = rl.file.Write(p)
	if err == nil {
		rl.currentSize += int64(n)
	}
	return n, err
}

func (rl *RotatingLogger) rotate() error {
	if err := rl.file.Close(); err != nil {
		return err
	}

	ext := filepath.Ext(rl.basePath)
	baseName := strings.TrimSuffix(rl.basePath, ext)
	timestamp := time.Now().Format("20060102_150405")

	for i := rl.backupCount - 1; i >= 0; i-- {
		var oldPath, newPath string
		
		if i == 0 {
			oldPath = rl.basePath
		} else {
			oldPath = fmt.Sprintf("%s.%d%s", baseName, i, ext)
		}
		
		if rl.compressOld && i == rl.backupCount-1 && fileExists(oldPath) {
			compressPath := oldPath + ".gz"
			if err := compressFile(oldPath, compressPath); err == nil {
				os.Remove(oldPath)
			}
			continue
		}

		if i == rl.backupCount-1 {
			if fileExists(oldPath) {
				os.Remove(oldPath)
			}
			continue
		}

		newPath = fmt.Sprintf("%s.%d%s", baseName, i+1, ext)
		if fileExists(oldPath) {
			os.Rename(oldPath, newPath)
		}
	}

	rotatedPath := fmt.Sprintf("%s.%s%s", baseName, timestamp, ext)
	os.Rename(rl.basePath, rotatedPath)

	file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	rl.file = file
	rl.currentSize = 0
	return nil
}

func (rl *RotatingLogger) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Close()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := NewGzipWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

type GzipWriter struct {
	io.WriteCloser
}

func NewGzipWriter(w io.Writer) *GzipWriter {
	return &GzipWriter{WriteCloser: &dummyWriter{w}}
}

type dummyWriter struct {
	io.Writer
}

func (d *dummyWriter) Close() error {
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log", 10, 5, true)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	log.SetOutput(logger)

	for i := 0; i < 1000; i++ {
		log.Printf("Log entry %d: %s", i, strings.Repeat("X", 1024))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation completed")
}package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingLogger struct {
	currentFile *os.File
	currentSize int64
	basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
	rl := &RotatingLogger{basePath: path}
	if err := rl.openCurrentFile(); err != nil {
		return nil, err
	}
	return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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
	oldPath := rl.basePath + "." + timestamp
	if err := os.Rename(rl.basePath, oldPath); err != nil {
		return err
	}

	if err := rl.compressFile(oldPath); err != nil {
		return err
	}

	if err := rl.cleanupOldFiles(); err != nil {
		return err
	}

	return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(src string) error {
	dst := src + ".gz"
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gz := gzip.NewWriter(dstFile)
	defer gz.Close()

	_, err = io.Copy(gz, srcFile)
	if err != nil {
		return err
	}

	return os.Remove(src)
}

func (rl *RotatingLogger) cleanupOldFiles() error {
	pattern := rl.basePath + ".*.gz"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	if len(matches) > maxBackups {
		toDelete := matches[:len(matches)-maxBackups]
		for _, f := range toDelete {
			os.Remove(f)
		}
	}
	return nil
}

func (rl *RotatingLogger) openCurrentFile() error {
	file, err := os.OpenFile(rl.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (rl *RotatingLogger) Close() error {
	if rl.currentFile != nil {
		return rl.currentFile.Close()
	}
	return nil
}

func main() {
	logger, err := NewRotatingLogger("app.log")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
		logger.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}