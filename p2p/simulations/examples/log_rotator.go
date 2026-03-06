
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type LogRotator struct {
    currentFile   *os.File
    currentSize   int64
    basePath      string
    currentNumber int
}

func NewLogRotator(basePath string) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath: basePath,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    rotator.findLatestBackupNumber()
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        err := lr.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    backupPath := lr.generateBackupPath()
    err := os.Rename(lr.basePath, backupPath)
    if err != nil {
        return err
    }

    err = lr.compressBackup(backupPath)
    if err != nil {
        return err
    }

    err = lr.cleanupOldBackups()
    if err != nil {
        fmt.Printf("Warning: cleanup failed: %v\n", err)
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = stat.Size()
    return nil
}

func (lr *LogRotator) generateBackupPath() string {
    timestamp := time.Now().Format("20060102_150405")
    lr.currentNumber++
    return fmt.Sprintf("%s.%s.%d", lr.basePath, timestamp, lr.currentNumber)
}

func (lr *LogRotator) compressBackup(backupPath string) error {
    src, err := os.Open(backupPath)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(backupPath + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    _, err = io.Copy(gz, src)
    if err != nil {
        return err
    }

    return os.Remove(backupPath)
}

func (lr *LogRotator) findLatestBackupNumber() {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil || len(matches) == 0 {
        lr.currentNumber = 0
        return
    }

    maxNum := 0
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }

        numStr := parts[len(parts)-2]
        num, err := strconv.Atoi(numStr)
        if err == nil && num > maxNum {
            maxNum = num
        }
    }

    lr.currentNumber = maxNum
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        err := os.Remove(matches[i])
        if err != nil {
            return err
        }
    }

    return nil
}

func (lr *LogRotator) Close() error {
    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("application.log")
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        _, err := rotator.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Write failed: %v\n", err)
            break
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}