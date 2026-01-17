package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize  = 10 * 1024 * 1024 // 10MB
    maxBackups   = 5
    currentLog   = "app.log"
    logPrefix    = "app"
    logExtension = ".log"
)

type LogRotator struct {
    currentSize int64
}

func NewLogRotator() *LogRotator {
    return &LogRotator{}
}

func (lr *LogRotator) Write(p []byte) (n int, err error) {
    if lr.currentSize+int64(len(p)) > maxFileSize {
        if err := lr.rotate(); err != nil {
            return 0, err
        }
    }

    file, err := os.OpenFile(currentLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    n, err = file.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if err := lr.cleanupOldLogs(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    newName := fmt.Sprintf("%s_%s%s", logPrefix, timestamp, logExtension)

    if err := os.Rename(currentLog, newName); err != nil {
        return err
    }

    lr.currentSize = 0
    return nil
}

func (lr *LogRotator) cleanupOldLogs() error {
    files, err := filepath.Glob(logPrefix + "_*" + logExtension)
    if err != nil {
        return err
    }

    sort.Sort(sort.Reverse(sort.StringSlice(files)))

    for i := maxBackups; i < len(files); i++ {
        if err := os.Remove(files[i]); err != nil {
            return err
        }
    }
    return nil
}

func (lr *LogRotator) loadCurrentSize() error {
    info, err := os.Stat(currentLog)
    if os.IsNotExist(err) {
        lr.currentSize = 0
        return nil
    }
    if err != nil {
        return err
    }
    lr.currentSize = info.Size()
    return nil
}

func main() {
    rotator := NewLogRotator()
    if err := rotator.loadCurrentSize(); err != nil {
        fmt.Printf("Failed to load current log size: %v\n", err)
        return
    }

    testMessage := fmt.Sprintf("[%s] Test log entry\n", time.Now().Format(time.RFC3339))
    _, err := rotator.Write([]byte(testMessage))
    if err != nil {
        fmt.Printf("Failed to write log: %v\n", err)
        return
    }

    fmt.Println("Log rotation test completed")
}package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const (
    maxFileSize = 10 * 1024 * 1024 // 10MB
    maxBackups  = 5
)

type RotatingLog struct {
    filePath   string
    current    *os.File
    currentSize int64
}

func NewRotatingLog(path string) (*RotatingLog, error) {
    rl := &RotatingLog{filePath: path}
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

    n, err := rl.current.Write(p)
    rl.currentSize += int64(n)
    return n, err
}

func (rl *RotatingLog) rotate() error {
    if err := rl.current.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.filePath, timestamp)
    if err := os.Rename(rl.filePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    if err := rl.cleanupOldBackups(); err != nil {
        return err
    }

    return rl.openCurrent()
}

func (rl *RotatingLog) compressFile(path string) error {
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

func (rl *RotatingLog) cleanupOldBackups() error {
    pattern := rl.filePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= maxBackups {
        return nil
    }

    var backupFiles []struct {
        path string
        time time.Time
    }

    for _, match := range matches {
        parts := strings.Split(match, ".")
        if len(parts) < 3 {
            continue
        }
        timestamp := parts[len(parts)-2]
        t, err := time.Parse("20060102_150405", timestamp)
        if err != nil {
            continue
        }
        backupFiles = append(backupFiles, struct {
            path string
            time time.Time
        }{match, t})
    }

    for i := 0; i < len(backupFiles)-maxBackups; i++ {
        if err := os.Remove(backupFiles[i].path); err != nil {
            return err
        }
    }

    return nil
}

func (rl *RotatingLog) openCurrent() error {
    file, err := os.OpenFile(rl.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    rl.current = file
    rl.currentSize = info.Size()
    return nil
}

func (rl *RotatingLog) Close() error {
    if rl.current != nil {
        return rl.current.Close()
    }
    return nil
}

func main() {
    log, err := NewRotatingLog("application.log")
    if err != nil {
        panic(err)
    }
    defer log.Close()

    for i := 0; i < 1000; i++ {
        message := fmt.Sprintf("Log entry %d: %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := log.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}