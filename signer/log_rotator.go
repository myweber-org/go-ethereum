
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    filePath    string
    maxSize     int64
    currentSize int64
    file        *os.File
}

func NewRotatingLogger(filePath string, maxSizeMB int) (*RotatingLogger, error) {
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

    return &RotatingLogger{
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
        file:        file,
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
    backupPath := rl.filePath + "." + timestamp

    if err := os.Rename(rl.filePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0

    go rl.cleanOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanOldLogs() {
    dir := filepath.Dir(rl.filePath)
    baseName := filepath.Base(rl.filePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    for _, entry := range entries {
        if !entry.IsDir() && len(entry.Name()) > len(baseName) {
            if entry.Name()[:len(baseName)] == baseName && entry.Name()[len(baseName)] == '.' {
                path := filepath.Join(dir, entry.Name())
                if info, err := os.Stat(path); err == nil {
                    if time.Since(info.ModTime()) > 7*24*time.Hour {
                        os.Remove(path)
                    }
                }
            }
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 1; i <= 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Test message for rotation check\n",
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}