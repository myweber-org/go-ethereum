package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type LogRotator struct {
    currentFile   *os.File
    filePath      string
    maxSize       int64
    rotationCount int
    maxRotations  int
}

func NewLogRotator(filePath string, maxSize int64, maxRotations int) (*LogRotator, error) {
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    return &LogRotator{
        currentFile:  file,
        filePath:     filePath,
        maxSize:      maxSize,
        maxRotations: maxRotations,
    }, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    if err := lr.checkRotation(); err != nil {
        return 0, err
    }
    return lr.currentFile.Write(p)
}

func (lr *LogRotator) checkRotation() error {
    info, err := lr.currentFile.Stat()
    if err != nil {
        return err
    }

    if info.Size() >= lr.maxSize {
        if err := lr.rotate(); err != nil {
            return err
        }
    }
    return nil
}

func (lr *LogRotator) rotate() error {
    if err := lr.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

    if err := os.Rename(lr.filePath, rotatedPath); err != nil {
        return err
    }

    if err := lr.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(lr.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    lr.currentFile = file
    lr.rotationCount++

    if lr.rotationCount > lr.maxRotations {
        if err := lr.cleanupOldFiles(); err != nil {
            return err
        }
    }

    return nil
}

func (lr *LogRotator) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, sourceFile); err != nil {
        return err
    }

    if err := os.Remove(sourcePath); err != nil {
        return err
    }

    return nil
}

func (lr *LogRotator) cleanupOldFiles() error {
    pattern := lr.filePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) > lr.maxRotations {
        filesToRemove := matches[:len(matches)-lr.maxRotations]
        for _, file := range filesToRemove {
            if err := os.Remove(file); err != nil {
                return err
            }
        }
    }

    return nil
}

func (lr *LogRotator) Close() error {
    return lr.currentFile.Close()
}

func main() {
    rotator, err := NewLogRotator("app.log", 1024*1024, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}