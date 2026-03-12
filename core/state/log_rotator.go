
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
    filename    string
    currentSize int64
    file        *os.File
}

func NewLogRotator(filename string) (*LogRotator, error) {
    rotator := &LogRotator{filename: filename}
    if err := rotator.openFile(); err != nil {
        return nil, err
    }
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
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

func (lr *LogRotator) openFile() error {
    info, err := os.Stat(lr.filename)
    if os.IsNotExist(err) {
        lr.file, err = os.Create(lr.filename)
        lr.currentSize = 0
        return err
    }
    if err != nil {
        return err
    }

    lr.currentSize = info.Size()
    lr.file, err = os.OpenFile(lr.filename, os.O_APPEND|os.O_WRONLY, 0644)
    return err
}

func (lr *LogRotator) rotate() error {
    if err := lr.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    rotatedFile := lr.filename + "." + timestamp

    if err := os.Rename(lr.filename, rotatedFile); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedFile); err != nil {
        return err
    }

    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }

    return lr.openFile()
}

func (lr *LogRotator) compressFile(source string) error {
    srcFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(source + ".gz")
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, srcFile)
    if err != nil {
        return err
    }

    if err := os.Remove(source); err != nil {
        return err
    }

    return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.filename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    var backupFiles []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }

        timestamp := parts[len(parts)-2]
        t, err := time.Parse("20060102150405", timestamp)
        if err != nil {
            continue
        }

        backupFiles = append(backupFiles, struct {
            path string
            time time.Time
        }{match, t})
    }

    for i := 0; i < len(backupFiles)-maxBackups; i++ {
        oldestIndex := 0
        for j := 1; j < len(backupFiles); j++ {
            if backupFiles[j].time.Before(backupFiles[oldestIndex].time) {
                oldestIndex = j
            }
        }

        if err := os.Remove(backupFiles[oldestIndex].path); err != nil {
            return err
        }

        backupFiles = append(backupFiles[:oldestIndex], backupFiles[oldestIndex+1:]...)
    }

    return nil
}

func (lr *LogRotator) Close() error {
    if lr.file != nil {
        return lr.file.Close()
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
        logEntry := fmt.Sprintf("[%s] Log entry %d: %s\n",
            time.Now().Format(time.RFC3339),
            i,
            strings.Repeat("X", 1024))

        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }

        if i%100 == 0 {
            fmt.Printf("Written %d log entries\n", i+1)
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}