package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentFile *os.File
    currentSize int64
    fileIndex   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath:  basePath,
        maxSize:   maxSize,
        fileIndex: 0,
    }

    err := logger.openCurrentFile()
    if err != nil {
        return nil, err
    }

    go logger.cleanupOldFiles()
    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        l.currentFile.Close()
    }

    filename := fmt.Sprintf("%s.log", l.basePath)
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.currentSize = info.Size()
    return nil
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentSize+int64(len(p)) > l.maxSize {
        err := l.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err := l.currentFile.Write(p)
    if err == nil {
        l.currentSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) rotate() error {
    if l.currentFile != nil {
        l.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedName := fmt.Sprintf("%s_%s.log", l.basePath, timestamp)
    err := os.Rename(fmt.Sprintf("%s.log", l.basePath), rotatedName)
    if err != nil {
        return err
    }

    err = l.compressFile(rotatedName)
    if err != nil {
        return err
    }

    l.fileIndex++
    return l.openCurrentFile()
}

func (l *RotatingLogger) compressFile(filename string) error {
    src, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(filename + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    _, err = io.Copy(gz, src)
    if err != nil {
        return err
    }

    os.Remove(filename)
    return nil
}

func (l *RotatingLogger) cleanupOldFiles() {
    for {
        time.Sleep(24 * time.Hour)

        files, err := filepath.Glob(l.basePath + "_*.log.gz")
        if err != nil {
            continue
        }

        cutoff := time.Now().Add(-30 * 24 * time.Hour)
        for _, file := range files {
            parts := strings.Split(file, "_")
            if len(parts) < 3 {
                continue
            }

            timestampStr := parts[len(parts)-2]
            t, err := time.Parse("20060102", timestampStr[:8])
            if err != nil {
                continue
            }

            if t.Before(cutoff) {
                os.Remove(file)
            }
        }
    }
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.currentFile != nil {
        return l.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(time.Millisecond * 100)
    }
}