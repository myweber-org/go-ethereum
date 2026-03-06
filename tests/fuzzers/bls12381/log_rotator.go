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
    file        *os.File
    basePath    string
    maxSize     int64
    currentSize int64
    sequence    int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
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
        sequence:    0,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedName := fmt.Sprintf("%s.%s.%d.gz", rl.basePath, timestamp, rl.sequence)
    rl.sequence++

    if err := compressFile(rl.basePath, archivedName); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil && !os.IsNotExist(err) {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    return nil
}

func compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
        os.Exit(1)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}
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
	mu           sync.Mutex
	currentSize  int64
	maxSize      int64
	basePath     string
	currentFile  *os.File
	fileCount    int
}

func NewRotatingWriter(basePath string, maxSize int64) (*RotatingWriter, error) {
	w := &RotatingWriter{
		basePath: basePath,
		maxSize:  maxSize,
	}
	if err := w.openCurrentFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotatingWriter) openCurrentFile() error {
	w.fileCount++
	filename := fmt.Sprintf("%s.%d.log", w.basePath, w.fileCount)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
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

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.currentSize+int64(len(p)) > w.maxSize && w.currentSize > 0 {
		w.currentFile.Close()
		if err := w.openCurrentFile(); err != nil {
			return 0, err
		}
	}

	n, err = w.currentFile.Write(p)
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
	writer, err := NewRotatingWriter("app", 1024*1024)
	if err != nil {
		fmt.Printf("Failed to create writer: %v\n", err)
		return
	}
	defer writer.Close()

	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
		writer.Write([]byte(msg))
		time.Sleep(10 * time.Millisecond)
	}
}