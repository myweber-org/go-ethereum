
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

const (
    maxFileSize    = 10 * 1024 * 1024 // 10MB
    maxBackupFiles = 5
    logDir         = "./logs"
)

type RotatingLogger struct {
    mu       sync.Mutex
    file     *os.File
    filePath string
    size     int64
}

func NewRotatingLogger(baseName string) (*RotatingLogger, error) {
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }

    filePath := filepath.Join(logDir, baseName+".log")
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }

    stat, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

    return &RotatingLogger{
        file:     file,
        filePath: filePath,
        size:     stat.Size(),
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

    backupFiles, err := rl.getBackupFiles()
    if err != nil {
        return err
    }

    for i := len(backupFiles) - 1; i >= 0; i-- {
        oldName := backupFiles[i]
        parts := strings.Split(oldName, ".")
        if len(parts) < 3 {
            continue
        }

        num, err := strconv.Atoi(parts[len(parts)-2])
        if err != nil {
            continue
        }

        if num >= maxBackupFiles-1 {
            os.Remove(filepath.Join(logDir, oldName))
        } else {
            newName := strings.Join(parts[:len(parts)-2], ".") + "." + strconv.Itoa(num+1) + ".log"
            os.Rename(
                filepath.Join(logDir, oldName),
                filepath.Join(logDir, newName),
            )
        }
    }

    newName := strings.TrimSuffix(rl.filePath, ".log") + ".1.log"
    if err := os.Rename(rl.filePath, newName); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.filePath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.file = file
    rl.size = 0
    return nil
}

func (rl *RotatingLogger) getBackupFiles() ([]string, error) {
    files, err := os.ReadDir(logDir)
    if err != nil {
        return nil, err
    }

    base := filepath.Base(rl.filePath)
    prefix := strings.TrimSuffix(base, ".log") + "."
    var backups []string

    for _, file := range files {
        if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) && strings.HasSuffix(file.Name(), ".log") {
            backups = append(backups, file.Name())
        }
    }

    sort.Slice(backups, func(i, j int) bool {
        numI := extractNumber(backups[i])
        numJ := extractNumber(backups[j])
        return numI < numJ
    })

    return backups, nil
}

func extractNumber(filename string) int {
    parts := strings.Split(filename, ".")
    if len(parts) < 3 {
        return 0
    }
    num, _ := strconv.Atoi(parts[len(parts)-2])
    return num
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    return rl.file.Close()
}

func main() {
    logger, err := NewRotatingLogger("app")
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 1000; i++ {
        msg := fmt.Sprintf("[%s] Log entry %d: Test message for rotation check\n",
            time.Now().Format("2006-01-02 15:04:05"), i)
        if _, err := logger.Write([]byte(msg)); err != nil {
            fmt.Printf("Write error: %v\n", err)
            break
        }
        time.Sleep(10 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}