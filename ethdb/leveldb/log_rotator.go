
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

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compressOld bool) (*LogRotator, error) {
    maxSizeBytes := int64(maxSizeMB) * 1024 * 1024

    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSizeBytes,
        maxBackups:  maxBackups,
        compressOld: compressOld,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
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

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.rotationMutex.Lock()
    defer lr.rotationMutex.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
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

    timestamp := time.Now().Format("20060102150405")
    backupPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    err := os.Rename(lr.basePath, backupPath)
    if err != nil {
        return err
    }

    err = lr.openCurrentFile()
    if err != nil {
        return err
    }

    if lr.compressOld {
        go lr.compressBackup(backupPath)
    }

    lr.cleanupOldBackups()
    return nil
}

func (lr *LogRotator) compressBackup(backupPath string) {
    compressedPath := backupPath + ".gz"

    src, err := os.Open(backupPath)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(compressedPath)
    if err != nil {
        return
    }
    defer dst.Close()

    gzWriter := gzip.NewWriter(dst)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, src)
    if err != nil {
        return
    }

    os.Remove(backupPath)
}

func (lr *LogRotator) cleanupOldBackups() {
    dir := filepath.Dir(lr.basePath)
    baseName := filepath.Base(lr.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") {
            backups = append(backups, name)
        }
    }

    if len(backups) <= lr.maxBackups {
        return
    }

    backupsToRemove := backups[:len(backups)-lr.maxBackups]
    for _, backup := range backupsToRemove {
        os.Remove(filepath.Join(dir, backup))
    }
}

func (lr *LogRotator) Close() error {
    lr.rotationMutex.Lock()
    defer lr.rotationMutex.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n",
            time.Now().Format(time.RFC3339), i)
        _, err := rotator.Write([]byte(logEntry))
        if err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }

        if i%100 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
    }

    fmt.Println("Log rotation demonstration completed")
}