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

type RotatingLogger struct {
    mu            sync.Mutex
    basePath      string
    maxSize       int64
    currentSize   int64
    currentFile   *os.File
    sequence      int
    compressOld   bool
}

func NewRotatingLogger(basePath string, maxSizeMB int, compressOld bool) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    logger := &RotatingLogger{
        basePath:    basePath,
        maxSize:     maxSize,
        compressOld: compressOld,
        sequence:    0,
    }

    err := logger.openOrCreateLog()
    return logger, err
}

func (rl *RotatingLogger) openOrCreateLog() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

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

    rl.findLatestSequence()
    return nil
}

func (rl *RotatingLogger) findLatestSequence() {
    pattern := rl.basePath + ".*"
    matches, _ := filepath.Glob(pattern)
    maxSeq := 0

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 2 {
            continue
        }

        seqStr := parts[len(parts)-1]
        if strings.HasSuffix(seqStr, ".gz") {
            seqStr = seqStr[:len(seqStr)-3]
        }

        seq, err := strconv.Atoi(seqStr)
        if err == nil && seq > maxSeq {
            maxSeq = seq
        }
    }

    rl.sequence = maxSeq
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
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

func (rl *RotatingLogger) rotate() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    rl.sequence++
    rotatedPath := fmt.Sprintf("%s.%d", rl.basePath, rl.sequence)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if rl.compressOld {
        go rl.compressFile(rotatedPath)
    }

    return rl.openOrCreateLog()
}

func (rl *RotatingLogger) compressFile(path string) {
    src, err := os.Open(path)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(path + ".gz")
    if err != nil {
        return
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return
    }

    os.Remove(path)
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app.log", 10, true)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Application event occurred\n",
            time.Now().Format("2006-01-02 15:04:05"), i)
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}