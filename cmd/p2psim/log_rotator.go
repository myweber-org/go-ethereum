
package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
    logFileName = "app.log"
)

type RotatingLogger struct {
    currentFile *os.File
    basePath    string
    fileSize    int64
}

func NewRotatingLogger(basePath string) (*RotatingLogger, error) {
    if err := os.MkdirAll(basePath, 0755); err != nil {
        return nil, err
    }

    fullPath := filepath.Join(basePath, logFileName)
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
        basePath:    basePath,
        fileSize:    info.Size(),
    }, nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    if rl.fileSize+int64(len(p)) > maxFileSize {
        if err := rl.rotate(); err != nil {
            log.Printf("Failed to rotate log: %v", err)
        }
    }

    n, err := rl.currentFile.Write(p)
    if err == nil {
        rl.fileSize += int64(n)
    }
    return n, err
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    backupPath := filepath.Join(rl.basePath, fmt.Sprintf("%s.%s", logFileName, timestamp))

    if err := os.Rename(filepath.Join(rl.basePath, logFileName), backupPath); err != nil {
        return err
    }

    file, err := os.OpenFile(filepath.Join(rl.basePath, logFileName), os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.fileSize = 0

    go rl.cleanupOldLogs()

    return nil
}

func (rl *RotatingLogger) cleanupOldLogs() {
    files, err := filepath.Glob(filepath.Join(rl.basePath, logFileName+".*"))
    if err != nil {
        return
    }

    if len(files) <= maxBackups {
        return
    }

    var timestamps []time.Time
    fileMap := make(map[time.Time]string)

    for _, file := range files {
        parts := strings.Split(file, ".")
        if len(parts) < 2 {
            continue
        }

        tsStr := parts[len(parts)-1]
        ts, err := time.Parse("20060102_150405", tsStr)
        if err != nil {
            continue
        }

        timestamps = append(timestamps, ts)
        fileMap[ts] = file
    }

    if len(timestamps) <= maxBackups {
        return
    }

    for i := 0; i < len(timestamps)-maxBackups; i++ {
        if filePath, exists := fileMap[timestamps[i]]; exists {
            os.Remove(filePath)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}

func main() {
    logger, err := NewRotatingLogger("./logs")
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()

    log.SetOutput(io.MultiWriter(os.Stdout, logger))

    for i := 0; i < 1000; i++ {
        log.Printf("Log entry %d: %s", i, strconv.FormatInt(time.Now().UnixNano(), 10))
        time.Sleep(10 * time.Millisecond)
    }
}