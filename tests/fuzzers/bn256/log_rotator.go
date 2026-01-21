
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
    mu           sync.Mutex
    currentFile  *os.File
    filePath     string
    maxSize      int64
    currentSize  int64
    rotationCount int
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
        currentFile:  file,
        filePath:     filePath,
        maxSize:      maxSize,
        currentSize:  info.Size(),
        rotationCount: 0,
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
    
    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s.%d", rl.filePath, timestamp, rl.rotationCount)
    
    if err := os.Rename(rl.filePath, backupPath); err != nil {
        return err
    }
    
    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    rl.currentFile = file
    rl.currentSize = 0
    rl.rotationCount++
    
    rl.cleanOldBackups()
    
    return nil
}

func (rl *RotatingLogger) cleanOldBackups() {
    dir := filepath.Dir(rl.filePath)
    baseName := filepath.Base(rl.filePath)
    
    files, err := os.ReadDir(dir)
    if err != nil {
        return
    }
    
    var backups []string
    for _, file := range files {
        if !file.IsDir() {
            name := file.Name()
            if len(name) > len(baseName) && name[:len(baseName)] == baseName {
                backups = append(backups, filepath.Join(dir, name))
            }
        }
    }
    
    if len(backups) > 10 {
        for i := 0; i < len(backups)-10; i++ {
            os.Remove(backups[i])
        }
    }
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
    logger, err := NewRotatingLogger("app.log", 10)
    if err != nil {
        panic(err)
    }
    defer logger.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry %d: Some sample log data here\n", 
            time.Now().Format(time.RFC3339), i)
        logger.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}