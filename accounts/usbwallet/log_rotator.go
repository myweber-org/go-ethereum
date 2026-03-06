
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
    current   *os.File
    size      int64
    mu        sync.Mutex
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        filename: filename,
    }
    
    if err := rl.openFile(); err != nil {
        return nil, err
    }
    
    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }
    
    n, err = rl.current.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if rl.current != nil {
        rl.current.Close()
    }
    
    timestamp := time.Now().Format("20060102-150405")
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
    
    return rl.openFile()
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
    
    if len(matches) > maxBackups {
        toRemove := matches[:len(matches)-maxBackups]
        for _, file := range toRemove {
            os.Remove(file)
        }
    }
    
    return nil
}

func (rl *RotatingLogger) openFile() error {
    f, err := os.OpenFile(rl.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    
    info, err := f.Stat()
    if err != nil {
        f.Close()
        return err
    }
    
    rl.current = f
    rl.size = info.Size()
    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}