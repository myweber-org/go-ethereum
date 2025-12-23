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
    currentSize int64
    currentFile *os.File
    fileCounter int
}

func NewRotatingLogger(basePath string, maxSizeMB int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath: basePath,
        maxSize:  maxSize,
    }

    err := logger.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return logger, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

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

    rl.scanExistingRotatedFiles()
    return nil
}

func (rl *RotatingLogger) scanExistingRotatedFiles() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    maxCounter := 0
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") {
            parts := strings.Split(name, ".")
            if len(parts) >= 3 && parts[len(parts)-1] == "gz" {
                if counter, err := strconv.Atoi(parts[len(parts)-2]); err == nil && counter > maxCounter {
                    maxCounter = counter
                }
            }
        }
    }
    rl.fileCounter = maxCounter
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    rl.fileCounter++
    rotatedPath := fmt.Sprintf("%s.%d.gz", rl.basePath, rl.fileCounter)

    oldFile, err := os.Open(rl.basePath)
    if err != nil {
        return err
    }
    defer oldFile.Close()

    rotatedFile, err := os.Create(rotatedPath)
    if err != nil {
        return err
    }
    defer rotatedFile.Close()

    gzWriter := gzip.NewWriter(rotatedFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, oldFile)
    if err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    return rl.openCurrentFile()
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
    logger, err := NewRotatingLogger("/var/log/myapp/app.log", 10)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
        os.Exit(1)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application event processed\n",
            time.Now().Format("2006-01-02 15:04:05"), i)
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}