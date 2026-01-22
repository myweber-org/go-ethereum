
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
    maxFiles    int
    currentFile *os.File
    currentSize int64
}

func NewRotatingLog(basePath string, maxSizeMB int, maxFiles int) (*RotatingLog, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLog{
        basePath:    basePath,
        maxSize:     maxSize,
        maxFiles:    maxFiles,
        currentFile: file,
        currentSize: info.Size(),
    }, nil
}

func (r *RotatingLog) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize+int64(len(p)) > r.maxSize {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.currentFile.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *RotatingLog) rotate() error {
    if err := r.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", r.basePath, timestamp)

    if err := os.Rename(r.basePath, rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(r.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    r.currentFile = file
    r.currentSize = 0

    go r.compressAndCleanup(rotatedPath)

    return nil
}

func (r *RotatingLog) compressAndCleanup(path string) {
    compressedPath := path + ".gz"

    source, err := os.Open(path)
    if err != nil {
        return
    }
    defer source.Close()

    dest, err := os.Create(compressedPath)
    if err != nil {
        return
    }
    defer dest.Close()

    gz := gzip.NewWriter(dest)
    defer gz.Close()

    if _, err := io.Copy(gz, source); err != nil {
        return
    }

    os.Remove(path)
    r.cleanupOldFiles()
}

func (r *RotatingLog) cleanupOldFiles() {
    dir := filepath.Dir(r.basePath)
    baseName := filepath.Base(r.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var compressedFiles []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            compressedFiles = append(compressedFiles, filepath.Join(dir, name))
        }
    }

    if len(compressedFiles) <= r.maxFiles {
        return
    }

    for i := 0; i < len(compressedFiles)-r.maxFiles; i++ {
        os.Remove(compressedFiles[i])
    }
}

func (r *RotatingLog) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.currentFile.Close()
}

func main() {
    log, err := NewRotatingLog("app.log", 10, 5)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        log.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}