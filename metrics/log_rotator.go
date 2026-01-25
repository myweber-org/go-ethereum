
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

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    mu          sync.Mutex
    currentSize int64
    file        *os.File
    basePath    string
    sequence    int
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
    rl := &RotatingLogger{
        basePath: basePath,
    }

    if err := rl.openCurrent(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrent() error {
    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.file = file
    rl.currentSize = info.Size()
    rl.sequence = rl.getLatestSequence()

    return nil
}

func (rl *RotatingLogger) getLatestSequence() int {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return 0
    }

    maxSeq := 0
    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        seq, err := strconv.Atoi(parts[len(parts)-2])
        if err == nil && seq > maxSeq {
            maxSeq = seq
        }
    }
    return maxSeq
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.currentSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102-150405")
    rl.sequence++
    backupPath := fmt.Sprintf("%s.%d.%s.gz", rl.basePath, rl.sequence, timestamp)

    if err := compressFile(rl.basePath, backupPath); err != nil {
        return err
    }

    if err := os.Remove(rl.basePath); err != nil {
        return err
    }

    if err := rl.cleanOldBackups(); err != nil {
        return err
    }

    return rl.openCurrent()
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

func (rl *RotatingLogger) cleanOldBackups() error {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    for i := 0; i < len(matches)-maxBackups; i++ {
        if err := os.Remove(matches[i]); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.file != nil {
        return rl.file.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("application.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", 
            time.Now().Format(time.RFC3339), i)
        if _, err := logger.Write([]byte(logEntry)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }
}