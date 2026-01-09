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
    maxBackups   = 5
    currentLog   = "app.log"
    logPrefix    = "app"
    logExtension = ".log"
)

type LogRotator struct {
    currentSize int64
}

func NewLogRotator() *LogRotator {
    return &LogRotator{}
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    file, err := os.OpenFile(currentLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
    if err := lr.cleanupOldLogs(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    newName := fmt.Sprintf("%s_%s%s", logPrefix, timestamp, logExtension)

    if err := os.Rename(currentLog, newName); err != nil {
        return err
    }

    lr.currentSize = 0
    return nil
}

func (lr *LogRotator) cleanupOldLogs() error {
    files, err := filepath.Glob(logPrefix + "_*" + logExtension)
    if err != nil {
        return err
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))

    for i := maxBackups; i < len(files); i++ {
        if err := os.Remove(files[i]); err != nil {
            return err
        }
    }
    return nil
}

func (lr *LogRotator) loadCurrentSize() error {
    info, err := os.Stat(currentLog)
    if os.IsNotExist(err) {
        lr.currentSize = 0
        return nil
    }
    if err != nil {
        return err
    }
    lr.currentSize = info.Size()
    return nil
}

func main() {
    rotator := NewLogRotator()
    if err := rotator.loadCurrentSize(); err != nil {
        fmt.Printf("Failed to load current log size: %v\n", err)
        return
    }

    testMessage := fmt.Sprintf("[%s] Test log entry\n", time.Now().Format(time.RFC3339))
    _, err := rotator.Write([]byte(testMessage))
    if err != nil {
        fmt.Printf("Failed to write log: %v\n", err)
        return
    }

    fmt.Println("Log rotation test completed")
}