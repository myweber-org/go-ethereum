
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
	maxBackups  = 5
)

type RotatingWriter struct {
	mu          sync.Mutex
	currentSize int64
	basePath    string
	currentFile *os.File
	backupCount int
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
	w := &RotatingWriter{
		basePath: path,
	}
	if err := w.openCurrent(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openCurrent() error {
	file, err := os.OpenFile(w.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}
	w.currentFile = file
	w.currentSize = stat.Size()
	return nil
}

func (w *RotatingWriter) rotate() error {
	w.currentFile.Close()

	for i := maxBackups - 1; i >= 0; i-- {
		oldPath := w.backupPath(i)
		newPath := w.backupPath(i + 1)

		if _, err := os.Stat(oldPath); err == nil {
			if i == maxBackups-1 {
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	os.Rename(w.basePath, w.backupPath(0))
	return w.openCurrent()
}

func (w *RotatingWriter) backupPath(index int) string {
	if index == 0 {
		return w.basePath + ".1"
	}
	return fmt.Sprintf("%s.%d.gz", w.basePath, index)
}

func (w *RotatingWriter) compressOldLogs() {
	for i := 1; i <= maxBackups; i++ {
		path := fmt.Sprintf("%s.%d", w.basePath, i)
		if _, err := os.Stat(path); err == nil {
			go compressFile(path)
		}
	}
}

func compressFile(path string) {
	// Compression implementation would go here
	// For simplicity, we'll just rename to .gz extension
	if !strings.HasSuffix(path, ".gz") {
		os.Rename(path, path+".gz")
	}
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > maxFileSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
		w.compressOldLogs()
	}

	n, err := w.currentFile.Write(p)
	if err == nil {
		w.currentSize += int64(n)
	}
	return n, err
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.currentFile.Close()
}

func main() {
	writer, err := NewRotatingWriter("app.log")
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(logEntry))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}