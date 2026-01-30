package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 1024 * 1024 // 1MB
    maxBackupFiles = 5
    logFileName   = "app.log"
)

type LogRotator struct {
    currentSize int64
    file        *os.File
}

func NewLogRotator() (*LogRotator, error) {
    file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        return nil, err
    }

    return &LogRotator{
        currentSize: info.Size(),
        file:        file,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxLogSize {
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
    lr.file.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    os.Rename(logFileName, backupName)

    file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    lr.file = file
    lr.currentSize = 0

    go lr.cleanupOldLogs()
    return nil
}

func (lr *LogRotator) cleanupOldLogs() {
    files, err := filepath.Glob(logFileName + ".*")
    if err != nil {
        return
    }

    if len(files) <= maxBackupFiles {
        return
    }

    for i := 0; i < len(files)-maxBackupFiles; i++ {
        os.Remove(files[i])
    }
}

func (lr *LogRotator) Close() error {
    return lr.file.Close()
}

func main() {
    rotator, err := NewLogRotator()
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d\n", time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}