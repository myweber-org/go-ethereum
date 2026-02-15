package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    filePath    string
    maxSize     int64
    maxAge      time.Duration
    currentFile *os.File
    currentSize int64
    mu          sync.Mutex
}

func NewRotator(filePath string, maxSize int64, maxAge time.Duration) (*Rotator, error) {
    r := &Rotator{
        filePath: filePath,
        maxSize:  maxSize,
        maxAge:   maxAge,
    }
    if err := r.openCurrent(); err != nil {
        return nil, err
    }
    go r.timeBasedRotation()
    return r, nil
}

func (r *Rotator) openCurrent() error {
    dir := filepath.Dir(r.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    f, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    stat, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    r.currentFile = f
    r.currentSize = stat.Size()
    return nil
}

func (r *Rotator) rotate() error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentFile != nil {
        r.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)
    if err := os.Rename(r.filePath, backupPath); err != nil {
        return err
    }

    return r.openCurrent()
}

func (r *Rotator) timeBasedRotation() {
    ticker := time.NewTicker(r.maxAge)
    defer ticker.Stop()
    for range ticker.C {
        if err := r.rotate(); err != nil {
            fmt.Fprintf(os.Stderr, "Rotation error: %v\n", err)
        }
    }
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.currentSize >= r.maxSize {
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

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewRotator("logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        panic(err)
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        msg := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}