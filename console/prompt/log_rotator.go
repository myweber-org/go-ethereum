package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackupCount = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentSize int64
    basePath    string
}

func NewLogRotator(basePath string) *LogRotator {
    return &LogRotator{
        basePath: basePath,
    }
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    file, err := os.OpenFile(filepath.Join(lr.basePath, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    n, err = file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    baseLogPath := filepath.Join(lr.basePath, logFileName)
    timestamp := time.Now().Format("20060102_150405")

    // Rename current log
    backupPath := filepath.Join(lr.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))
    if err := os.Rename(baseLogPath, backupPath); err != nil && !os.IsNotExist(err) {
        return err
    }

    // Reset current size
    lr.currentSize = 0

    // Cleanup old backups
    return lr.cleanupOldBackups()
}

func (lr *LogRotator) cleanupOldBackups() error {
    files, err := filepath.Glob(filepath.Join(lr.basePath, logFileName+".*"))
    if err != nil {
        return err
    }

    // Sort by timestamp (newest first)
    sort.Slice(files, func(i, j int) bool {
        return extractTimestamp(files[i]) > extractTimestamp(files[j])
    })

    // Remove excess backups
    for i := maxBackupCount; i < len(files); i++ {
        if err := os.Remove(files[i]); err != nil {
            return err
        }
    }
    return nil
}

func extractTimestamp(path string) string {
    parts := strings.Split(filepath.Base(path), ".")
    if len(parts) > 1 {
        return parts[len(parts)-1]
    }
    return ""
}

func main() {
    rotator := NewLogRotator(".")
    
    // Simulate log writing
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Iteration %d: Sample log message\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        
        // Simulate time passing
        if i%100 == 0 {
            time.Sleep(10 * time.Millisecond)
        }
    }
    
    fmt.Println("Log rotation test completed")
}