
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxFileSize = 1024 * 1024 // 1MB
	maxBackups  = 5
)

type RotatingWriter struct {
	mu       sync.Mutex
	filename string
	file     *os.File
	size     int64
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
	w := &RotatingWriter{filename: filename}
	if err := w.openFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openFile() error {
	info, err := os.Stat(w.filename)
	if os.IsNotExist(err) {
		w.size = 0
	} else if err != nil {
		return err
	} else {
		w.size = info.Size()
	}

	file, err := os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.size+int64(len(p)) > maxFileSize {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *RotatingWriter) rotate() error {
	if w.file != nil {
		w.file.Close()
	}

	for i := maxBackups - 1; i >= 0; i-- {
		oldName := w.backupName(i)
		newName := w.backupName(i + 1)

		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(w.filename, w.backupName(0)); err != nil && !os.IsNotExist(err) {
		return err
	}

	w.size = 0
	return w.openFile()
}

func (w *RotatingWriter) backupName(i int) string {
	if i == 0 {
		return w.filename
	}
	ext := filepath.Ext(w.filename)
	base := w.filename[:len(w.filename)-len(ext)]
	return fmt.Sprintf("%s.%d%s", base, i, ext)
}

func (w *RotatingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func main() {
	writer, err := NewRotatingWriter("app.log")
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("[%s] Log entry %d: This is a sample log message.\n",
			time.Now().Format(time.RFC3339), i)
		if _, err := writer.Write([]byte(msg)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}