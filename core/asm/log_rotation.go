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
    maxFileSize = 1024 * 1024 // 1MB
    maxFiles    = 5
    logDir      = "./logs"
)

type RotatingLogger struct {
    currentFile *os.File
    baseName    string
    fileSize    int64
    fileIndex   int
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    rl := &RotatingLogger{
        baseName: baseName,
    }

    if err := rl.openNextFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openNextFile() error {
    if rl.currentFile != nil {
        rl.currentFile.Close()
    }

    rl.fileIndex = (rl.fileIndex + 1) % maxFiles
    fileName := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, rl.fileIndex))

    file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.fileSize = 0

    rl.cleanupOldFiles()

    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() {
    for i := 0; i < maxFiles; i++ {
        if i == rl.fileIndex {
            continue
        }
        fileName := filepath.Join(logDir, fmt.Sprintf("%s_%d.log", rl.baseName, i))
        os.Remove(fileName)
    }
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
    if rl.fileSize+int64(len(p)) > maxFileSize {
        if err := rl.openNextFile(); err != nil {
            return 0, err
        }
    }

    n, err = rl.currentFile.Write(p)
    rl.fileSize += int64(n)
    return n, err
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    multiWriter := io.MultiWriter(os.Stdout, logger)
    log.SetOutput(multiWriter)

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d at %s", i, time.Now().Format(time.RFC3339))
        time.Sleep(10 * time.Millisecond)
    }
}