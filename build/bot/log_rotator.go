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
    file        *os.File
    maxSize     int64
    basePath    string
    currentSize int64
    fileCount   int
}

func NewRotatingWriter(basePath string, maxSize int64) (*RotatingWriter, error) {
    w := &RotatingWriter{
        maxSize:  maxSize,
        basePath: basePath,
    }
    if err := w.openFile(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) openFile() error {
    w.fileCount++
    filename := fmt.Sprintf("%s.%d", w.basePath, w.fileCount)
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    w.file = file
    stat, err := file.Stat()
    if err != nil {
        return err
    }
    w.currentSize = stat.Size()
    return nil
}

func (w *RotatingWriter) rotate() error {
    if w.file != nil {
        w.file.Close()
    }
    return w.openFile()
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    if w.currentSize+int64(len(p)) > w.maxSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = w.file.Write(p)
    w.currentSize += int64(n)
    return n, err
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
    writer, err := NewRotatingWriter("app.log", 1024*1024)
    if err != nil {
        panic(err)
    }
    defer writer.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}