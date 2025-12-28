package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

const (
    maxFileSize = 1024 * 1024 // 1MB
    maxBackups  = 5
)

type RotatingLogger struct {
    filename   string
    current    *os.File
    size       int64
}

func NewRotatingLogger(filename string) (*RotatingLogger, error) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        return nil, err
    }

    return &RotatingLogger{
        filename: filename,
        current:  file,
        size:     info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.size+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.current.Write(p)
    rl.size += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.current.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", rl.filename, timestamp)
    if err := os.Rename(rl.filename, backupName); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filename, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.current = file
    rl.size = 0

    go rl.cleanupOldBackups()

    return nil
}

func (rl *RotatingLogger) cleanupOldBackups() {
    pattern := rl.filename + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackups {
        return
    }

    backups := make([]backupInfo, 0, len(matches))
    for _, match := range matches {
        info, err := os.Stat(match)
        if err != nil {
            continue
        }
        backups = append(backups, backupInfo{
            path: match,
            time: info.ModTime(),
        })
    }

    for i := 0; i < len(backups)-maxBackups; i++ {
        os.Remove(backups[i].path)
    }
}

type backupInfo struct {
    path string
    time time.Time
}

func (rl *RotatingLogger) Close() error {
    return rl.current.Close()
}

func main() {
    logger, err := NewRotatingLogger("app.log")
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    for i := 1; i <= 10000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        logger.Write([]byte(message))
        time.Sleep(10 * time.Millisecond)
    }
}