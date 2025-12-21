
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
    
    lr.cleanupOldBackups()
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
    
    backupsToRemove := len(backups) - lr.maxBackups
    for i := 0; i < backupsToRemove; i++ {
        path := filepath.Join(dir, backups[i])
        os.Remove(path)
    }
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
    rotator, err := NewLogRotator("app.log", 10, 5, true)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()
    
    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        
        if i%100 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
    }
    
    fmt.Println("Log rotation test completed")
}