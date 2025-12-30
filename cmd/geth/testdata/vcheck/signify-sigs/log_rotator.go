
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotatingWriter struct {
	mu          sync.Mutex
	currentFile *os.File
	filePath    string
	maxSize     int64
	currentSize int64
	fileCount   int
	maxFiles    int
}

func NewRotatingWriter(basePath string, maxSize int64, maxFiles int) (*RotatingWriter, error) {
	writer := &RotatingWriter{
		filePath: basePath,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	if err := writer.openCurrentFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

func (w *RotatingWriter) openCurrentFile() error {
	file, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	w.currentFile = file
	w.currentSize = info.Size()
	return nil
}

func (w *RotatingWriter) rotate() error {
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	for i := w.maxFiles - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", w.filePath, i)
		newPath := fmt.Sprintf("%s.%d", w.filePath, i+1)

		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	backupPath := fmt.Sprintf("%s.1", w.filePath)
	os.Rename(w.filePath, backupPath)

	return w.openCurrentFile()
}

func (w *RotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
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

	if w.currentFile != nil {
		return w.currentFile.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("app.log", 1024*1024, 5)
	if err != nil {
		fmt.Printf("Failed to create rotating writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(logEntry))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Log rotation test completed")
}