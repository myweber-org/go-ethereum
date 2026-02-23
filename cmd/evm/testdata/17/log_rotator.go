
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLogger struct {
    mu         sync.Mutex
    file       *os.File
    currentDir string
    baseName   string
    size       int64
}

func NewRotatingLogger(dir, name string) (*RotatingLogger, error) {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    basePath := filepath.Join(dir, name)
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
        file:       file,
        currentDir: dir,
        baseName:   name,
        size:       info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.file.Write(p)
    if err == nil {
        rl.size += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.file.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Unix()
    rotatedName := fmt.Sprintf("%s.%d", rl.baseName, timestamp)
    rotatedPath := filepath.Join(rl.currentDir, rotatedName)

    if err := os.Rename(filepath.Join(rl.currentDir, rl.baseName), rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(filepath.Join(rl.currentDir, rl.baseName), os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.size = 0
    rl.cleanupOldBackups()
    return nil
}

func (rl *RotatingLogger) compressFile(path string) error {
    src, err := os.Open(path)
    if err != nil {
        return err
    }
    defer src.Close()

    dst, err := os.Create(path + ".gz")
    if err != nil {
        return err
    }
    defer dst.Close()

    gz := gzip.NewWriter(dst)
    defer gz.Close()

    if _, err := io.Copy(gz, src); err != nil {
        return err
    }

    if err := os.Remove(path); err != nil {
        return err
    }

    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := filepath.Join(rl.currentDir, rl.baseName+".*.gz")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    var backups []struct {
        path string
        time int64
    }

    for _, match := range matches {
        base := filepath.Base(match)
        suffix := base[len(rl.baseName)+1 : len(base)-3]
        timestamp, err := strconv.ParseInt(suffix, 10, 64)
        if err != nil {
            continue
        }
        backups = append(backups, struct {
            path string
            time int64
        }{match, timestamp})
    }

    for i := 0; i < len(backups)-maxBackups; i++ {
        oldestIdx := i
        for j := i + 1; j < len(backups); j++ {
            if backups[j].time < backups[oldestIdx].time {
                oldestIdx = j
            }
        }
        backups[i], backups[oldestIdx] = backups[oldestIdx], backups[i]
        os.Remove(backups[i].path)
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("./logs", "app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(msg))
        time.Sleep(10 * time.Millisecond)
    }
}