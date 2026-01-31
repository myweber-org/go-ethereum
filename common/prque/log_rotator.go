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

const maxFileSize = 10 * 1024 * 1024 // 10MB
const maxBackupFiles = 5

type LogRotator struct {
    filePath    string
    currentFile *os.File
    currentSize int64
}

func NewLogRotator(path string) (*LogRotator, error) {
    rotator := &LogRotator{filePath: path}
    if err := rotator.openCurrentFile(); err != nil {
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

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

    if err := os.Rename(lr.filePath, backupPath); err != nil {
        return err
    }

    if err := lr.compressBackup(backupPath); err != nil {
        return err
    }

    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressBackup(src string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destFile, err := os.Create(src + ".gz")
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    os.Remove(src)
    return nil
}

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.filePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackupFiles {
        return nil
    }

    var backupFiles []string
    for _, match := range matches {
        backupFiles = append(backupFiles, match)
    }

    sortByTimestamp(backupFiles)

    for i := 0; i < len(backupFiles)-maxBackupFiles; i++ {
        os.Remove(backupFiles[i])
    }

    return nil
}

func sortByTimestamp(files []string) {
    for i := 0; i < len(files); i++ {
        for j := i + 1; j < len(files); j++ {
            ts1 := extractTimestamp(files[i])
            ts2 := extractTimestamp(files[j])
            if ts1 > ts2 {
                files[i], files[j] = files[j], files[i]
            }
        }
    }
}

func extractTimestamp(filename string) int64 {
    base := filepath.Base(filename)
    parts := strings.Split(base, ".")
    if len(parts) < 3 {
        return 0
    }

    tsStr := parts[1]
    ts, err := strconv.ParseInt(tsStr, 10, 64)
    if err != nil {
        return 0
    }
    return ts
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
    lr.currentSize = info.Size()
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
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}