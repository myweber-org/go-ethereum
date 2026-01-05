package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type LogRotator struct {
    filePath    string
    maxSize     int64
    maxBackups  int
    currentSize int64
    file        *os.File
    mu          sync.Mutex
}

func NewLogRotator(filePath string, maxSizeMB int, maxBackups int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &LogRotator{
        filePath:    filePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := lr.file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

    if err := os.Rename(lr.filePath, backupPath); err != nil {
        return err
    }

    if err := lr.compressBackup(backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.file = file
    lr.currentSize = 0

    lr.cleanupOldBackups()

    return nil
}

func (lr *LogRotator) compressBackup(srcPath string) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destPath := srcPath + ".gz"
    destFile, err := os.Create(destPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    if err := os.Remove(srcPath); err != nil {
        return err
    }

    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := lr.filePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= lr.maxBackups {
        return
    }

    backups := make([]struct {
        path string
        time time.Time
    }, 0, len(matches))

    for _, match := range matches {
        base := filepath.Base(match)
        suffix := base[len(filepath.Base(lr.filePath))+1 : len(base)-3]
        if t, err := time.Parse("20060102150405", suffix); err == nil {
            backups = append(backups, struct {
                path string
                time time.Time
            }{path: match, time: t})
        }
    }

    for i := 0; i < len(backups)-lr.maxBackups; i++ {
        os.Remove(backups[i].path)
    }
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator("app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 1; i <= 1000; i++ {
        logEntry := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }

        if i%100 == 0 {
            fmt.Printf("Written %d log entries\n", i)
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}