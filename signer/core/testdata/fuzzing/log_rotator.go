package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    logDirectory = "./logs"
)

type LogRotator struct {
    currentFile *os.File
    filePath    string
    bytesWritten int64
}

func NewLogRotator(baseName string) (*LogRotator, error) {
    if err := os.MkdirAll(logDirectory, 0755); err != nil {
        return nil, err
    }
    
    filePath := filepath.Join(logDirectory, baseName+".log")
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    info, _ := file.Stat()
    return &LogRotator{
        currentFile: file,
        filePath:    filePath,
        bytesWritten: info.Size(),
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.bytesWritten+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.bytesWritten += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    lr.currentFile.Close()
    
    timestamp := time.Now().Format("20060102-150405")
    rotatedPath := filepath.Join(logDirectory, fmt.Sprintf("%s.%s.log", 
        filepath.Base(lr.filePath[:len(lr.filePath)-4]), timestamp))
    
    if err := os.Rename(lr.filePath, rotatedPath); err != nil {
        return err
    }
    
    file, err := os.OpenFile(lr.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    lr.currentFile = file
    lr.bytesWritten = 0
    
    go lr.cleanupOldFiles()
    return nil
}

func (lr *LogRotator) cleanupOldFiles() {
    pattern := filepath.Join(logDirectory, filepath.Base(lr.filePath[:len(lr.filePath)-4])+".*.log")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }
    
    if len(matches) > maxBackups {
        filesToDelete := matches[:len(matches)-maxBackups]
        for _, file := range filesToDelete {
            os.Remove(file)
        }
    }
}

func (lr *LogRotator) Close() error {
    return lr.currentFile.Close()
}

func main() {
    rotator, err := NewLogRotator("application")
    if err != nil {
        panic(err)
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Test message for rotation\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}