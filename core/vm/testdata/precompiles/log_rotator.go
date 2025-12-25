
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

type RotatingLogger struct {
    mu           sync.Mutex
    file         *os.File
    currentSize  int64
    maxSize      int64
    basePath     string
    rotationCount int
    maxRotations int
}

func NewRotatingLogger(basePath string, maxSize int64, maxRotations int) (*RotatingLogger, error) {
    file, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        file:         file,
        currentSize:  info.Size(),
        maxSize:      maxSize,
        basePath:     basePath,
        maxRotations: maxRotations,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentSize+int64(len(p)) > rl.maxSize {
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

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    compressedPath := rotatedPath + ".gz"
    if err := compressFile(rotatedPath, compressedPath); err != nil {
        return err
    }

    os.Remove(rotatedPath)

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.currentSize = 0
    rl.rotationCount++

    if rl.rotationCount > rl.maxRotations {
        rl.cleanOldRotations()
    }

    return nil
}

func compressFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()

    gz := gzip.NewWriter(destination)
    defer gz.Close()

    _, err = io.Copy(gz, source)
    return err
}

func (rl *RotatingLogger) cleanOldRotations() {
    pattern := rl.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > rl.maxRotations {
        filesToRemove := matches[:len(matches)-rl.maxRotations]
        for _, file := range filesToRemove {
            os.Remove(file)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log", 1024*1024, 5)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        logger.Write([]byte(fmt.Sprintf("Log entry %d: Application event occurred at %v\n", i, time.Now())))
        time.Sleep(10 * time.Millisecond)
    }
}