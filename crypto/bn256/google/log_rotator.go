
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

type LogRotator struct {
    basePath      string
    maxSize       int64
    maxBackups    int
    compressOld   bool
    currentSize   int64
    currentFile   *os.File
    rotationMutex sync.Mutex
}

func NewLogRotator(basePath string, maxSize int64, maxBackups int, compressOld bool) (*LogRotator, error) {
    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compressOld,
    }

    if err := rotator.openCurrentFile(); err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.rotationMutex.Lock()
    defer lr.rotationMutex.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
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

    timestamp := time.Now().Format("2006-01-02_15-04-05")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    if err := os.Rename(lr.basePath, rotatedPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    if lr.compressOld {
        if err := lr.compressFile(rotatedPath); err != nil {
            return fmt.Errorf("failed to compress log file: %w", err)
        }
        rotatedPath = rotatedPath + ".gz"
    }

    if err := lr.cleanupOldBackups(rotatedPath); err != nil {
        return fmt.Errorf("failed to cleanup old backups: %w", err)
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    destFile, err := os.Create(sourcePath + ".gz")
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    if err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldBackups(newestBackup string) error {
    dir := filepath.Dir(lr.basePath)
    baseName := filepath.Base(lr.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && name != filepath.Base(newestBackup) {
            backups = append(backups, filepath.Join(dir, name))
        }
    }

    if len(backups) > lr.maxBackups {
        sortBackups(backups)
        for i := 0; i < len(backups)-lr.maxBackups; i++ {
            os.Remove(backups[i])
        }
    }

    return nil
}

func sortBackups(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(path string) int64 {
    base := filepath.Base(path)
    parts := strings.Split(base, ".")
    if len(parts) < 2 {
        return 0
    }

    timestampStr := parts[len(parts)-1]
    if strings.HasSuffix(timestampStr, ".gz") {
        timestampStr = timestampStr[:len(timestampStr)-3]
    }

    timestampStr = strings.ReplaceAll(timestampStr, "-", "")
    timestampStr = strings.ReplaceAll(timestampStr, "_", "")

    ts, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }
    return ts
}

func (lr *LogRotator) openCurrentFile() error {
    file, err := os.OpenFile(lr.basePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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
    lr.rotationMutex.Lock()
    defer lr.rotationMutex.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}