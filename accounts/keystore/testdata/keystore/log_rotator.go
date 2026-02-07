
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
    mu          sync.Mutex
    file        *os.File
    filePath    string
    maxSize     int64
    currentSize int64
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
        file:        file,
        filePath:    filePath,
        maxSize:     maxSize,
        currentSize: info.Size(),
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
    backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
    
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
    
    files, err := os.ReadDir(dir)
    if err != nil {
        return
    }
    
    for _, file := range files {
        if file.IsDir() {
            continue
        }
        
        name := file.Name()
        if len(name) > len(baseName) && name[:len(baseName)] == baseName && name[len(baseName)] == '.' {
            filePath := filepath.Join(dir, name)
            info, err := os.Stat(filePath)
            if err != nil {
                continue
            }
            
            if time.Since(info.ModTime()) > 30*24*time.Hour {
                os.Remove(filePath)
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
        panic(err)
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}