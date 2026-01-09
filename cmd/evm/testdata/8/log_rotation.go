package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

const (
    maxLogSize    = 1024 * 1024 // 1MB
    maxBackupLogs = 5
    logFileName   = "app.log"
)

type RotatingLogger struct {
    currentFile *os.File
    fileSize    int64
    basePath    string
}

func NewRotatingLogger(path string) (*RotatingLogger, error) {
    fullPath := filepath.Join(path, logFileName)
    file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        currentFile: file,
        fileSize:    info.Size(),
        basePath:    path,
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.fileSize+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.fileSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.currentFile.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupName := fmt.Sprintf("%s.%s", logFileName, timestamp)
    backupPath := filepath.Join(rl.basePath, backupName)
    originalPath := filepath.Join(rl.basePath, logFileName)

    if err := os.Rename(originalPath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(originalPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.fileSize = 0

    go rl.cleanupOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    pattern := filepath.Join(rl.basePath, logFileName+".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) <= maxBackupLogs {
        return
    }

    for i := 0; i < len(matches)-maxBackupLogs; i++ {
        os.Remove(matches[i])
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger(".")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(10 * time.Millisecond)
    }
}