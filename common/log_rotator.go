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

type RotatingLog struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    currentSize int64
    currentFile *os.File
    fileCount   int
}

func NewRotatingLog(basePath string, maxSizeMB int) (*RotatingLog, error) {
    rl := &RotatingLog{
        basePath: basePath,
        maxSize:  int64(maxSizeMB) * 1024 * 1024,
    }

    dir := filepath.Dir(basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLog) openCurrentFile() error {
    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.currentFile = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLog) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", rl.basePath, timestamp)

    oldFile, err := os.Open(rl.basePath)
    if err != nil {
        return err
    }
    defer oldFile.Close()

    archivedFile, err := os.Create(archivedPath)
    if err != nil {
        return err
    }
    defer archivedFile.Close()

    gzWriter := gzip.NewWriter(archivedFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, oldFile); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    rl.fileCount++
    return rl.openCurrentFile()
}

func (rl *RotatingLog) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func (rl *RotatingLog) CleanOldFiles(maxFiles int) error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var archivedFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            archivedFiles = append(archivedFiles, filepath.Join(dir, name))
        }
    }

    if len(archivedFiles) <= maxFiles {
        return nil
    }

    for i := 0; i < len(archivedFiles)-maxFiles; i++ {
        if err := os.Remove(archivedFiles[i]); err != nil {
            return err
        }
    }

    return nil
}

func main() {
    log, err := NewRotatingLog("/var/log/myapp/app.log", 10)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
            time.Now().Format(time.RFC3339), i)
        if _, err := log.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }

        if i%100 == 0 {
            if err := log.CleanOldFiles(5); err != nil {
                fmt.Printf("Clean error: %v\n", err)
            }
        }

        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}