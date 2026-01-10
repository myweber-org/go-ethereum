
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