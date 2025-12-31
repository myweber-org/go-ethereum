
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
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logFileName    = "app.log"
)

type LogRotator struct {
    currentFile *os.File
    filePath    string
    bytesWritten int64
}

func NewLogRotator() (*LogRotator, error) {
    lr := &LogRotator{
        filePath: logFileName,
    }
    
    if err := lr.openCurrentFile(); err != nil {
        return nil, err
    }
    
    return lr, nil
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }
    
    lr.currentFile = file
    lr.bytesWritten = info.Size()
    
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.bytesWritten+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := lr.currentFile.Write(p)
    if err != nil {
        return n, err
    }
    
    lr.bytesWritten += int64(n)
    return n, nil
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return err
    }
    
    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", lr.filePath, timestamp)
    
    if err := os.Rename(lr.filePath, backupName); err != nil {
        return err
    }
    
    if err := lr.openCurrentFile(); err != nil {
        return err
    }
    
    lr.cleanupOldBackups()
    
    return nil
}

func (lr *LogRotator) cleanupOldBackups() {
    pattern := fmt.Sprintf("%s.*", lr.filePath)
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) <= maxBackupFiles {
        return
    }
    
    sort.Sort(sort.Reverse(sort.StringSlice(matches)))
    
    for i := maxBackupFiles; i < len(matches); i++ {
        os.Remove(matches[i])
    }
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func extractTimestamp(filename string) (time.Time, error) {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return time.Time{}, fmt.Errorf("invalid backup filename")
    }
    
    timestampStr := parts[len(parts)-1]
    return time.Parse("20060102_150405", timestampStr)
}

func main() {
    rotator, err := NewLogRotator()
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        os.Exit(1)
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}