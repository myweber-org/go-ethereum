
package main

import (
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

type RotatingWriter struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    fileIndex   int
}

func NewRotatingWriter(path string) (*RotatingWriter, error) {
    writer := &RotatingWriter{
        basePath: path,
    }
    if err := writer.openCurrentFile(); err != nil {
        return nil, err
    }
    return writer, nil
}

func (w *RotatingWriter) openCurrentFile() error {
    if w.currentFile != nil {
        w.currentFile.Close()
    }

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
            if err := os.Rename(oldPath, newPath); err != nil {
                return err
            }
        }
    }

    if err := os.Rename(w.basePath, w.backupPath(0)); err != nil && !os.IsNotExist(err) {
        return err
    }

    return w.openCurrentFile()
}

func (w *RotatingWriter) backupPath(index int) string {
    if index == 0 {
        return w.basePath + ".1"
    }
    return fmt.Sprintf("%s.%d", w.basePath, index+1)
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
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
    if w.currentFile != nil {
        return w.currentFile.Close()
    }
    return nil
}

func main() {
    writer, err := NewRotatingWriter("app.log")
    if err != nil {
        fmt.Printf("Failed to create writer: %v\n", err)
        return
    }
    defer writer.Close()

    for i := 0; i < 100; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: This is a sample log message\n",
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(logEntry))
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}