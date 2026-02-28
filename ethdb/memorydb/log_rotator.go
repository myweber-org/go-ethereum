
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

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    filename   string
    current    *os.File
    size       int64
    mu         sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{filename: filename}
    if err := rl.openCurrent(); err != nil {
        return nil, err
    }
    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.current.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) openCurrent() error {
    info, err := os.Stat(rl.filename)
    if err != nil && !os.IsNotExist(err) {
        return err
    }

    if err == nil {
        rl.size = info.Size()
    } else {
        rl.size = 0
    }

    file, err := os.OpenFile(rl.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    rl.current = file
    return nil
}

func (rl *RotatingLogger) rotate() error {
    if rl.current != nil {
        rl.current.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", rl.filename, timestamp)

    if err := os.Rename(rl.filename, backupName); err != nil {
        return err
    }

    if err := rl.compressBackup(backupName); err != nil {
        return err
    }

    if err := rl.cleanOldBackups(); err != nil {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLogger) compressBackup(filename string) error {
    src, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(filename + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    os.Remove(filename)
    return nil
}

func (rl *RotatingLogger) cleanOldBackups() error {
    pattern := rl.filename + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        os.Remove(matches[i])
    }
    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}package main

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
    go r.ageMonitor()
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

func (r *Rotator) Write(p []byte) (int, error) {
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

func (r *Rotator) rotate() error {
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

func (r *Rotator) ageMonitor() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()
    for range ticker.C {
        r.mu.Lock()
        if r.currentFile != nil {
            stat, err := r.currentFile.Stat()
            if err == nil && time.Since(stat.ModTime()) > r.maxAge {
                r.rotate()
            }
        }
        r.mu.Unlock()
    }
}

func (r *Rotator) Close() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if r.currentFile != nil {
        return r.currentFile.Close()
    }
    return nil
}