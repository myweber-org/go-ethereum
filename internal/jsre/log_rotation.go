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
    if err := os.MkdirAll(filepath.Dir(basePath), 0755); err != nil {
        return nil, err
    }

    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    if err := logger.openCurrentFile(); err != nil {
        return nil, err
    }

    return logger, nil
}

func (l *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(l.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    l.currentFile = file
    l.fileSize = info.Size()
    return nil
}

func (l *RotatingLogger) rotate() error {
    l.currentFile.Close()

    timestamp := time.Now().Format("20060102_150405")
    archivePath := fmt.Sprintf("%s.%s.gz", l.basePath, timestamp)

    oldFile, err := os.Open(l.basePath)
    if err != nil {
        return err
    }
    defer oldFile.Close()

    archiveFile, err := os.Create(archivePath)
    if err != nil {
        return err
    }
    defer archiveFile.Close()

    gzWriter := gzip.NewWriter(archiveFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, oldFile); err != nil {
        return err
    }

    if err := os.Remove(l.basePath); err != nil {
        return err
    }

    l.fileCount++
    if l.fileCount > l.maxFiles {
        l.cleanupOldFiles()
    }

    return l.openCurrentFile()
}

func (l *RotatingLogger) cleanupOldFiles() {
    pattern := l.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > l.maxFiles {
        filesToDelete := matches[:len(matches)-l.maxFiles]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (l *RotatingLogger) Write(p []byte) (int, error) {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.fileSize+int64(len(p)) > l.maxSize {
        if err := l.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := l.currentFile.Write(p)
    if err == nil {
        l.fileSize += int64(n)
    }
    return n, err
}

func (l *RotatingLogger) Close() error {
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10*1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: Application event occurred at %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(100 * time.Millisecond)
    }
}