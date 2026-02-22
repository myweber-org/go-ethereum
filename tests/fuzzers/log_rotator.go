
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

const maxFileSize = 1024 * 1024 * 10 // 10MB

type LogRotator struct {
    basePath   string
    maxBackups int
}

func NewLogRotator(path string, backups int) *LogRotator {
    return &LogRotator{
        basePath:   path,
        maxBackups: backups,
    }
}

func (lr *LogRotator) Write(data []byte) error {
    currentPath := lr.basePath

    if shouldRotate(currentPath) {
        err := lr.rotate()
        if err != nil {
            return fmt.Errorf("rotation failed: %w", err)
        }
    }

    file, err := os.OpenFile(currentPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("open file failed: %w", err)
    }
    defer file.Close()

    _, err = file.Write(data)
    return err
}

func shouldRotate(path string) bool {
    info, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    if err != nil {
        return false
    }
    return info.Size() >= maxFileSize
}

func (lr *LogRotator) rotate() error {
    for i := lr.maxBackups - 1; i >= 0; i-- {
        oldPath := lr.backupPath(i)
        newPath := lr.backupPath(i + 1)

        if i == lr.maxBackups-1 {
            os.Remove(newPath)
            continue
        }

        if _, err := os.Stat(oldPath); err == nil {
            err := os.Rename(oldPath, newPath)
            if err != nil {
                return fmt.Errorf("rename failed: %w", err)
            }
        }
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)
    return os.Rename(lr.basePath, rotatedPath)
}

func (lr *LogRotator) backupPath(index int) string {
    if index == 0 {
        return lr.basePath
    }
    return fmt.Sprintf("%s.%d", lr.basePath, index)
}

func (lr *LogRotator) CleanOldBackups() error {
    pattern := lr.basePath + ".*"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    var backupFiles []string
    for _, match := range matches {
        if isTimestampBackup(match) {
            backupFiles = append(backupFiles, match)
        }
    }

    if len(backupFiles) > lr.maxBackups {
        toRemove := backupFiles[lr.maxBackups:]
        for _, file := range toRemove {
            os.Remove(file)
        }
    }
    return nil
}

func isTimestampBackup(path string) bool {
    parts := strings.Split(path, ".")
    if len(parts) < 2 {
        return false
    }

    lastPart := parts[len(parts)-1]
    if len(lastPart) != 15 {
        return false
    }

    if _, err := strconv.Atoi(lastPart[:8]); err != nil {
        return false
    }

    if _, err := strconv.Atoi(lastPart[9:]); err != nil {
        return false
    }

    return lastPart[8] == '_'
}

func main() {
    rotator := NewLogRotator("/var/log/app.log", 5)

    testData := []byte(fmt.Sprintf("Test log entry at %s\n", time.Now().Format(time.RFC3339)))
    err := rotator.Write(testData)
    if err != nil {
        fmt.Printf("Write error: %v\n", err)
    }

    err = rotator.CleanOldBackups()
    if err != nil {
        fmt.Printf("Cleanup error: %v\n", err)
    }

    fmt.Println("Log rotation completed")
}
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

type LogRotator struct {
    mu          sync.Mutex
    basePath    string
    maxSize     int64
    maxFiles    int
    currentSize int64
    currentFile *os.File
}

func NewLogRotator(basePath string, maxSizeMB int, maxFiles int) (*LogRotator, error) {
    maxSize := int64(maxSizeMB) * 1024 * 1024
    rotator := &LogRotator{
        basePath: basePath,
        maxSize:  maxSize,
        maxFiles: maxFiles,
    }

    err := rotator.openCurrentFile()
    if err != nil {
        return nil, err
    }

    return rotator, nil
}

func (lr *LogRotator) openCurrentFile() error {
    dir := filepath.Dir(lr.basePath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return err
    }

    file, err := os.OpenFile(lr.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    info, err := file.Stat()
    if err != nil {
        file.Close()
        return err
    }

    lr.currentFile = file
    lr.currentSize = info.Size()
    return nil
}

func (lr *LogRotator) Write(p []byte) (int, error) {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentSize+int64(len(p)) > lr.maxSize {
        err := lr.rotate()
        if err != nil {
            return 0, err
        }
    }

    n, err := lr.currentFile.Write(p)
    if err == nil {
        lr.currentSize += int64(n)
    }
    return n, err
}

func (lr *LogRotator) rotate() error {
    if lr.currentFile != nil {
        lr.currentFile.Close()
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", lr.basePath, timestamp)

    err := os.Rename(lr.basePath, rotatedPath)
    if err != nil {
        return err
    }

    err = lr.compressFile(rotatedPath)
    if err != nil {
        return err
    }

    err = lr.cleanupOldFiles()
    if err != nil {
        return err
    }

    return lr.openCurrentFile()
}

func (lr *LogRotator) compressFile(sourcePath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    compressedPath := sourcePath + ".gz"
    compressedFile, err := os.Create(compressedPath)
    if err != nil {
        return err
    }
    defer compressedFile.Close()

    gzWriter := gzip.NewWriter(compressedFile)
    defer gzWriter.Close()

    _, err = io.Copy(gzWriter, sourceFile)
    if err != nil {
        return err
    }

    os.Remove(sourcePath)
    return nil
}

func (lr *LogRotator) cleanupOldFiles() error {
    pattern := lr.basePath + ".*.gz"
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return err
    }

    if len(matches) <= lr.maxFiles {
        return nil
    }

    var timestamps []time.Time
    timestampMap := make(map[time.Time]string)

    for _, match := range matches {
        base := filepath.Base(match)
        tsStr := strings.TrimSuffix(strings.TrimPrefix(base, filepath.Base(lr.basePath)+"."), ".gz")
        ts, err := time.Parse("20060102_150405", tsStr)
        if err != nil {
            continue
        }
        timestamps = append(timestamps, ts)
        timestampMap[ts] = match
    }

    for i := 0; i < len(timestamps)-lr.maxFiles; i++ {
        oldest := timestamps[i]
        fileToRemove := timestampMap[oldest]
        os.Remove(fileToRemove)
    }

    return nil
}

func (lr *LogRotator) Close() error {
    lr.mu.Lock()
    defer lr.mu.Unlock()

    if lr.currentFile != nil {
        return lr.currentFile.Close()
    }
    return nil
}

func main() {
    rotator, err := NewLogRotator("/var/log/myapp/app.log", 10, 5)
    if err != nil {
        fmt.Printf("Failed to create log rotator: %v\n", err)
        return
    }
    defer rotator.Close()

    for i := 0; i < 1000; i++ {
        logEntry := fmt.Sprintf("[%s] Log entry number %d\n", time.Now().Format(time.RFC3339), i)
        rotator.Write([]byte(logEntry))
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}