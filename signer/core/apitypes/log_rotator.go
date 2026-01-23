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
    backupCount = 5
)

type RotatingWriter struct {
    filename    string
    currentFile *os.File
    currentSize int64
}

func NewRotatingWriter(filename string) (*RotatingWriter, error) {
    w := &RotatingWriter{filename: filename}
    if err := w.openCurrentFile(); err != nil {
        return nil, err
    }
    return w, nil
}

func (w *RotatingWriter) Write(p []byte) (n int, err error) {
    if w.currentSize+int64(len(p)) > maxFileSize {
        if err := w.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = w.currentFile.Write(p)
    w.currentSize += int64(n)
    return n, err
}

func (w *RotatingWriter) rotate() error {
    if w.currentFile != nil {
        w.currentFile.Close()
    }

    for i := backupCount - 1; i >= 0; i-- {
        oldName := w.backupFilename(i)
        newName := w.backupFilename(i + 1)

        if _, err := os.Stat(oldName); err == nil {
            os.Rename(oldName, newName)
        }
    }

    if err := os.Rename(w.filename, w.backupFilename(0)); err != nil && !os.IsNotExist(err) {
        return err
    }

    return w.openCurrentFile()
}

func (w *RotatingWriter) backupFilename(index int) string {
    if index == 0 {
        return w.filename + ".1"
    }
    return fmt.Sprintf("%s.%d", w.filename, index+1)
}

func (w *RotatingWriter) openCurrentFile() error {
    file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Application is running normally\n",
            time.Now().Format(time.RFC3339), i)
        writer.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}