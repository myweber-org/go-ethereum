
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type LogRotator struct {
    mu          sync.Mutex
    currentFile *os.File
    filePath    string
    maxSize     int64
    maxBackups  int
}

func NewLogRotator(filePath string, maxSize int64, maxBackups int) (*LogRotator, error) {
    rotator := &LogRotator{
        filePath:   filePath,
        maxSize:    maxSize,
        maxBackups: maxBackups,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentFile == nil {
        if err := lr.openCurrentFile(); err != nil {
            return 0, err
        }
    }

    fileInfo, err := lr.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if fileInfo.Size()+int64(len(p)) > lr.maxSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    return lr.currentFile.Write(p)
}

func (lr *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(lr.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(lr.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    lr.currentFile = file
    return nil
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        if err := lr.currentFile.Close(); err != nil {
            return err
        }
    }

    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", lr.filePath, timestamp)

    if err := os.Rename(lr.filePath, backupPath); err != nil {
        return err
    }

    if err := lr.compressFile(backupPath); err != nil {
        return err
    }

    if err := lr.cleanupOldBackups(); err != nil {
        return err
    }

    return lr.openCurrentFile()
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

func (lr *LogRotator) cleanupOldBackups() error {
    pattern := lr.filePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= lr.maxBackups {
        return nil
    }

    backupsToDelete := matches[:len(matches)-lr.maxBackups]
    for _, backup := range backupsToDelete {
        if err := os.Remove(backup); err != nil {
            return err
        }
    }

    return nil
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}