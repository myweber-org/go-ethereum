
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
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    maxBackups    int
    currentSize   int64
    currentFile   *os.File
    compressOld   bool
}

func NewLogRotator(basePath string, maxSizeMB int, maxBackups int, compress bool) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    
    rotator := &LogRotator{
        basePath:    basePath,
        maxSize:     maxSize,
        maxBackups:  maxBackups,
        compressOld: compress,
    }
    
    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }
    
    return rotator, nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    
    if lr.currentSize+int64(len(p)) > lr.maxSize {
        err := lr.rotate()
        if err != nil {
            return 0, err
        }
    }
    
    n, err := lr.currentFile.Write(p)
    if err != nil {
        return n, err
    }
    
    lr.currentSize += int64(n)
    return n, nil
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

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }
    
    timestamp := time.Now().Format("20060102150405")
    backupPath := lr.basePath + "." + timestamp
    
    err := os.Rename(lr.basePath, backupPath)
    if err != nil {
        return err
    }
    
    err = lr.openCurrentFile()
    if err != nil {
        return err
    }
    
    go lr.manageBackups(backupPath)
    
    return nil
}

func (lr *LogRotator) manageBackups(backupPath string) {
    if lr.compressOld {
        compressedPath := backupPath + ".gz"
        err := compressFile(backupPath, compressedPath)
        if err == nil {
            os.Remove(backupPath)
            backupPath = compressedPath
        }
    }
    
    backups := lr.listBackups()
    if len(backups) > lr.maxBackups {
        toRemove := backups[lr.maxBackups:]
        for _, backup := range toRemove {
            os.Remove(backup)
        }
    }
}

func (lr *LogRotator) listBackups() []string {
    pattern := lr.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return []string{}
    }
    
    var backups []string
    for _, match := range matches {
        if strings.HasSuffix(match, ".gz") || isTimestampBackup(match) {
            backups = append(backups, match)
        }
    }
    
    return backups
}

func isTimestampBackup(path string) bool {
    parts := strings.Split(path, ".")
    if len(parts) < 2 {
        return false
    }
    
    timestamp := parts[len(parts)-1]
    if len(timestamp) != 14 {
        return false
    }
    
    _, err := strconv.Atoi(timestamp)
    return err == nil
}

func compressFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    
    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()
    
    gzWriter := gzip.NewWriter(dstFile)
    defer gzWriter.Close()
    
    _, err = io.Copy(gzWriter, srcFile)
    return err
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()
    
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
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }
    
    fmt.Println("Log rotation test completed")
}