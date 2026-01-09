package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLog struct {
    currentFile *os.File
    currentSize int64
    basePath    string
    sequence    int
}

func NewRotatingLog(basePath string) (*RotatingLog, error) {
    rl := &RotatingLog{
        basePath: basePath,
        sequence: 0,
    }

    if err := rl.openCurrent(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLog) Write(p []byte) (int, error) {
    if rl.currentSize+int64(len(p)) > maxFileSize {
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
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    oldPath := rl.currentFile.Name()
    timestamp := time.Now().Format("20060102_150405")
    archivedPath := fmt.Sprintf("%s.%s.gz", oldPath, timestamp)

    if err := compressFile(oldPath, archivedPath); err != nil {
        return err
    }

    if err := os.Remove(oldPath); err != nil {
        return err
    }

    rl.sequence++
    if rl.sequence > maxBackups {
        rl.cleanupOldBackups()
    }

    return rl.openCurrent()
}

func compressFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    gz := gzip.NewWriter(out)
    defer gz.Close()

    _, err = io.Copy(gz, in)
    return err
}

func (rl *RotatingLog) openCurrent() error {
    path := rl.basePath
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func (rl *RotatingLog) cleanupOldBackups() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackups {
        toRemove := matches[:len(matches)-maxBackups]
        for _, path := range toRemove {
            os.Remove(path)
        }
    }
}

func (rl *RotatingLog) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("application.log")
    if err != nil {
        fmt.Printf("Failed to create log: %v\n", err)
        return
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }
}