
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
    currentSize int64
    fileCount   int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    if maxSize <= 0 {
        return nil, fmt.Errorf("maxSize must be positive")
    }

    rl := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLogger) rotate() error {
    rl.currentFile.Close()

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%d.%s.gz", rl.basePath, rl.fileCount, timestamp)
    rl.fileCount++

    if err := compressFile(rl.basePath, archivedPath); err != nil {
        return fmt.Errorf("compression failed: %v", err)
    }

    if err := os.Remove(rl.basePath); err != nil {
        return fmt.Errorf("failed to remove old log: %v", err)
    }

    return rl.openCurrentFile()
}

func compressFile(source, target string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
        rl.currentSize = 0
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    tempDir := os.TempDir()
    logPath := filepath.Join(tempDir, "application.log")

    logger, err := NewRotatingLogger(logPath, 1)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Sample log data for rotation testing\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write failed: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}