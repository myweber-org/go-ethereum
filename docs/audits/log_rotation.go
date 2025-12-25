package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    filePath   string
    maxSize    int64
    currentSize int64
    maxFiles   int
}

func NewRotatingLogger(filePath string, maxSize int64, maxFiles int) (*RotatingLogger, error) {
    if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
        return nil, err
    }

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        file:       file,
        filePath:   filePath,
        maxSize:    maxSize,
        currentSize: info.Size(),
        maxFiles:   maxFiles,
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
    rl.file.Close()

    for i := rl.maxFiles - 1; i > 0; i-- {
        oldPath := fmt.Sprintf("%s.%d", rl.filePath, i)
        newPath := fmt.Sprintf("%s.%d", rl.filePath, i+1)

        if _, err := os.Stat(oldPath); err == nil {
            os.Rename(oldPath, newPath)
        }
    }

    backupPath := fmt.Sprintf("%s.1", rl.filePath)
    os.Rename(rl.filePath, backupPath)

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("logs/app.log", 1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}