
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
    mu          sync.Mutex
    basePath    string
    currentFile *os.File
    maxSize     int64
    fileCount   int
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    if maxSizeMB <= 0 {
        maxSizeMB = 10
    }
    if maxFiles <= 0 {
        maxFiles = 5
    }

    rl := &RotatingLogger{
        basePath: basePath,
        maxSize:  int64(maxSizeMB) * 1024 * 1024,
        maxFiles: maxFiles,
    }

    if err := rl.openCurrentFile(); err != nil {
        return nil, err
    }

    return rl, nil
}

func (rl *RotatingLogger) openCurrentFile() error {
    dir := filepath.Dir(rl.basePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.OpenFile(rl.basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    rl.currentFile = file
    return nil
}

func (rl *RotatingLogger) Write(p []byte) (int, error) {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    info, err := rl.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if info.Size()+int64(len(p)) > rl.maxSize {
        if err := rl.rotate(); err != nil {
            return 0, err
        }
    }

    return rl.currentFile.Write(p)
}

func (rl *RotatingLogger) rotate() error {
    if err := rl.currentFile.Close(); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102_150405")
    rotatedPath := fmt.Sprintf("%s.%s", rl.basePath, timestamp)

    if err := os.Rename(rl.basePath, rotatedPath); err != nil {
        return err
    }

    if err := rl.compressFile(rotatedPath); err != nil {
        return err
    }

    rl.fileCount++
    if rl.fileCount > rl.maxFiles {
        if err := rl.cleanupOldFiles(); err != nil {
            return err
        }
    }

    return rl.openCurrentFile()
}

func (rl *RotatingLogger) compressFile(srcPath string) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    destPath := srcPath + ".gz"
    destFile, err := os.Create(destPath)
    if err != nil {
        return err
    }
    defer destFile.Close()

    gzWriter := gzip.NewWriter(destFile)
    defer gzWriter.Close()

    if _, err := io.Copy(gzWriter, srcFile); err != nil {
        return err
    }

    if err := os.Remove(srcPath); err != nil {
        return err
    }

    return nil
}

func (rl *RotatingLogger) cleanupOldFiles() error {
    dir := filepath.Dir(rl.basePath)
    baseName := filepath.Base(rl.basePath)

    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var compressedFiles []string
    for _, entry := range entries {
        if strings.HasPrefix(entry.Name(), baseName+".") && strings.HasSuffix(entry.Name(), ".gz") {
            compressedFiles = append(compressedFiles, entry.Name())
        }
    }

    if len(compressedFiles) > rl.maxFiles {
        filesToRemove := compressedFiles[:len(compressedFiles)-rl.maxFiles]
        for _, file := range filesToRemove {
            if err := os.Remove(filepath.Join(dir, file)); err != nil {
                return err
            }
        }
    }

    return nil
}

func (rl *RotatingLogger) parseFileNumber(filename string) int {
    parts := strings.Split(filename, ".")
    if len(parts) < 2 {
        return 0
    }

    numPart := parts[len(parts)-2]
    if num, err := strconv.Atoi(numPart); err == nil {
        return num
    }
    return 0
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
    logger, err := NewRotatingLogger("/var/log/myapp/app.log", 5, 3)
    if err != nil {
        fmt.Printf("Failed to create logger: %v\n", err)
        return
    }
    defer logger.Close()

    for i := 0; i < 100; i++ {
        message := fmt.Sprintf("Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
        if _, err := logger.Write([]byte(message)); err != nil {
            fmt.Printf("Write error: %v\n", err)
        }
        time.Sleep(100 * time.Millisecond)
    }

    fmt.Println("Log rotation test completed")
}