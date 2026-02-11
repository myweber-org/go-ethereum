
package main

import (
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type RotatingLogger struct {
    mu          sync.Mutex
    currentFile *os.File
    basePath    string
    maxSize     int64
    fileCount   int
    maxFiles    int
}

func NewRotatingLogger(basePath string, maxSizeMB int, maxFiles int) (*RotatingLogger, error) {
    if maxFiles < 1 {
        maxFiles = 10
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

    stat, err := rl.currentFile.Stat()
    if err != nil {
        return 0, err
    }

    if stat.Size()+int64(len(p)) > rl.maxSize {
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

    compressedPath := rotatedPath + ".gz"
    if err := compressFile(rotatedPath, compressedPath); err != nil {
        return err
    }
    os.Remove(rotatedPath)

    rl.fileCount++
    if rl.fileCount > rl.maxFiles {
        rl.cleanupOldFiles()
    }

    return rl.openCurrentFile()
}

func compressFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    gz := gzip.NewWriter(out)
    defer gz.Close()

    _, err = io.Copy(gz, in)
    return err
}

func (rl *RotatingLogger) cleanupOldFiles() {
    pattern := rl.basePath + ".*.gz"
    files, err := filepath.Glob(pattern)
    if err != nil {
        return
    }

    if len(files) > rl.maxFiles {
        filesToDelete := files[:len(files)-rl.maxFiles]
        for _, f := range filesToDelete {
            os.Remove(f)
        }
    }
}

func (rl *RotatingLogger) Close() error {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if rl.currentFile != nil {
        return rl.currentFile.Close()
    }
    return nil
}