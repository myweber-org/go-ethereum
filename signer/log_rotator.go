package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Rotator struct {
    filePath     string
    maxSize      int64
    rotationTime time.Duration
    currentSize  int64
    lastRotation time.Time
    file         *os.File
    mu           sync.Mutex
}

func NewRotator(filePath string, maxSize int64, rotationTime time.Duration) (*Rotator, error) {
    dir := filepath.Dir(filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create directory: %w", err)
    }

    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, fmt.Errorf("failed to stat log file: %w", err)
    }

    return &Rotator{
        filePath:     filePath,
        maxSize:      maxSize,
        rotationTime: rotationTime,
        currentSize:  info.Size(),
        lastRotation: time.Now(),
        file:         file,
    }, nil
}

func (r *Rotator) Write(p []byte) (int, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if r.shouldRotate() {
        if err := r.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := r.file.Write(p)
    if err == nil {
        r.currentSize += int64(n)
    }
    return n, err
}

func (r *Rotator) shouldRotate() bool {
    if r.currentSize >= r.maxSize {
        return true
    }
    if time.Since(r.lastRotation) >= r.rotationTime {
        return true
    }
    return false
}

func (r *Rotator) rotate() error {
    if err := r.file.Close(); err != nil {
        return fmt.Errorf("failed to close current file: %w", err)
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", r.filePath, timestamp)

    if err := os.Rename(r.filePath, backupPath); err != nil {
        return fmt.Errorf("failed to rename log file: %w", err)
    }

    file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to create new log file: %w", err)
    }

    r.file = file
    r.currentSize = 0
    r.lastRotation = time.Now()
    return nil
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.file.Close()
}

func main() {
    rotator, err := NewRotator("logs/app.log", 1024*1024, 24*time.Hour)
    if err != nil {
        fmt.Printf("Failed to create rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := rotator.Write([]byte(message)); err != nil {
            fmt.Printf("Failed to write log: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}