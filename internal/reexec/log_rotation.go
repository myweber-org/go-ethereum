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
    filePath    string
    logger      *log.Logger
    written     int64
}

func NewRotatingLogger(dir string) (*RotatingLogger, error) {
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, err
    }

    filePath := filepath.Join(dir, logFileName)
    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    rl := &RotatingLogger{
        currentFile: file,
        filePath:    filePath,
        written:     info.Size(),
        logger:      log.New(file, "", log.LstdFlags),
    }

    return rl, nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.written+int64(len(p)) > maxLogSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    n, err = rl.currentFile.Write(p)
    rl.written += int64(n)
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    rl.currentFile.Close()

    timestamp := time.Now().Format("20060102_150405")
    backupPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
    if err := os.Rename(rl.filePath, backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.logger.SetOutput(file)
    rl.written = 0

    go rl.cleanupOldLogs()
    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    dir := filepath.Dir(rl.filePath)
    pattern := filepath.Join(dir, logFileName+".*")

    matches, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(matches) > maxBackupLogs {
        toDelete := matches[:len(matches)-maxBackupLogs]
        for _, path := range toDelete {
            os.Remove(path)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    return rl.currentFile.Close()
}

func main() {
    logger, err := NewRotatingLogger("./logs")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 10000; i++ {
        log.Printf("Log entry %d: Application is running normally", i)
        time.Sleep(10 * time.Millisecond)
    }
}