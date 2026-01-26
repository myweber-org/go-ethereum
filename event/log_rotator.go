
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
    mu           sync.Mutex
    basePath     string
    currentFile  *os.File
    maxSize      int64
    currentSize  int64
    maxBackups   int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxBackups int) (*RotatingLogger, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024

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
        basePath:    basePath,
        currentFile: file,
        maxSize:     maxSize,
        currentSize: info.Size(),
        maxBackups:  maxBackups,
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

    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := filepath.Join(dir, fmt.Sprintf("%s.%s", baseName, timestamp))

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    rl.currentSize = 0

    go rl.compressOldLog(rotatedPath)
    go rl.cleanupOldBackups()

    return nil
}

func (rl *RotatingLogger) compressOldLog(path string) {
    compressedPath := path + ".gz"

    src, err := os.Open(path)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(compressedPath)
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

func (rl *RotatingLogger) cleanupOldBackups() {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return
    }

    var backups []string
    for _, entry := range entries {
        name := entry.Name()
        if strings.HasPrefix(name, baseName+".") && strings.HasSuffix(name, ".gz") {
            backups = append(backups, filepath.Join(dir, name))
        }
    }

    if len(backups) <= rl.maxBackups {
        return
    }

    sortBackupsByTimestamp(backups)

    for i := 0; i < len(backups)-rl.maxBackups; i++ {
        os.Remove(backups[i])
    }
}

func sortBackupsByTimestamp(backups []string) {
    for i := 0; i < len(backups); i++ {
        for j := i + 1; j < len(backups); j++ {
            if extractTimestamp(backups[i]) > extractTimestamp(backups[j]) {
                backups[i], backups[j] = backups[j], backups[i]
            }
        }
    }
}

func extractTimestamp(path string) int64 {
    base := filepath.Base(path)
    parts := strings.Split(base, ".")
    if len(parts) < 3 {
        return 0
    }

    timestampStr := parts[1]
    timestampStr = strings.ReplaceAll(timestampStr, "_", "")
    ts, err := strconv.ParseInt(timestampStr, 10, 64)
    if err != nil {
        return 0
    }
    return ts
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
    logger, err := NewRotatingLogger("app.log", 10, 5)
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